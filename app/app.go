package app

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"

	"github.com/google/go-tika/tika"
	"github.com/jovandeginste/docmgmt/bayes"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
	"github.com/thoas/go-funk"
)

type App struct {
	Configuration Configuration
	Classifier    *bayes.Classifier
	TikaServer    *tika.Server
}

func (a *App) LoadConfiguration(cfgFile string) error {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Search config in home directory with name ".cobra-example" (without extension).
		viper.AddConfigPath("$HOME/.docmgmt")
		viper.SetConfigName("config")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			panic(err)
		}
	}
	err := viper.Unmarshal(&a.Configuration)
	if err != nil {
		return err
	}

	err = a.LoadClassifier()
	if err != nil {
		return err
	}

	p, err := homedir.Expand("~/.docmgmt/tika-server-" + a.Configuration.TikaVersion + ".jar")
	if err != nil {
		return err
	}

	s, err := tika.NewServer(p, "")
	if err != nil {
		return err
	}

	a.TikaServer = s

	return nil
}

func (a *App) LoadClassifier() error {
	cp, err := homedir.Expand(a.Configuration.ClassifierData)
	if err != nil {
		return err
	}

	c, err := bayes.NewClassifierFromFile(cp)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
		c = bayes.NewClassifier()
	}

	a.Classifier = c
	return nil
}

func (a *App) SaveClassifier() error {
	cp, err := homedir.Expand(a.Configuration.ClassifierData)
	if err != nil {
		return err
	}

	return a.Classifier.WriteToFile(cp)
}

func (a *App) StartServer() error {
	if a.Configuration.Verbose {
		log.Println("Starting tika server version: " + a.Configuration.TikaVersion)
	}
	return a.TikaServer.Start(context.Background())
}

func (a *App) StopServer() {
	if a.Configuration.Verbose {
		log.Println("Stopping tika server")
	}
	a.TikaServer.Stop()
}

func (a *App) WriteMetadata(file string) error {
	metaPath, err := a.MetadataDir(file, true)
	if err != nil {
		return err
	}

	meta, body, err := a.Parse(file)
	if err != nil {
		return err
	}

	metaJSON := meta.JSON()

	ioutil.WriteFile(path.Join(metaPath, "metadata.json"), metaJSON, 0600)
	ioutil.WriteFile(path.Join(metaPath, "body.txt"), []byte(body), 0600)
	return nil
}

func (a *App) MetadataDir(file string, create bool) (string, error) {
	sum := sha256.Sum256([]byte(file))
	relPath := fmt.Sprintf("%x", sum)
	shard := funk.Shard(relPath, 2, 2, false)
	fullPath := path.Join(append([]string{a.Configuration.MetadataRoot}, shard...)...)

	src, err := os.Stat(fullPath)
	if err != nil {
		if !os.IsNotExist(err) {
			return fullPath, err
		}
		os.MkdirAll(fullPath, 0700)
		src, err = os.Stat(fullPath)
	}
	if !src.IsDir() {
		return fullPath, errors.New("metadatadir '" + fullPath + "' exists but is not a directory")
	}

	return fullPath, nil
}

func (a *App) ReadFileInfo(file string) (*Info, error) {
	metaPath, err := a.MetadataDir(file, false)
	if err != nil {
		return nil, err
	}

	return a.ReadFileInfoMetaPath(metaPath)
}

func (a *App) ReadFileInfoMetaPath(metaPath string) (i *Info, err error) {
	body, err := a.ReadFileBodyMetaPath(metaPath)
	if err != nil {
		return
	}

	metadata, err := a.ReadFileMetadata(metaPath)
	if err != nil {
		return
	}

	tags, err := a.ReadFileTagsMetaPath(metaPath)

	i = &Info{
		Body:     body,
		Metadata: metadata,
		Tags:     tags,
	}

	return
}

func (a *App) ReadFileBody(file string) (string, error) {
	metaPath, err := a.MetadataDir(file, false)
	if err != nil {
		return "", err
	}

	return a.ReadFileBodyMetaPath(metaPath)
}

func (a *App) ReadFileBodyMetaPath(metaPath string) (string, error) {
	body, err := ioutil.ReadFile(path.Join(metaPath, "body.txt"))
	if os.IsNotExist(err) {
		return "", errors.New("file not parsed yet")
	}

	return string(body), err
}

func (a *App) ReadFileMetadata(metaPath string) (*Metadata, error) {
	var m *Metadata

	meta, err := ioutil.ReadFile(path.Join(metaPath, "metadata.json"))
	if os.IsNotExist(err) {
		return m, errors.New("file not parsed yet")
	}

	err = json.Unmarshal(meta, &m)

	return m, err
}

func (a *App) ReadFileTags(file string) ([]string, error) {
	metaPath, err := a.MetadataDir(file, false)
	if err != nil {
		return []string{}, err
	}

	return a.ReadFileTagsMetaPath(metaPath)
}

func (a *App) ReadFileTagsMetaPath(metaPath string) ([]string, error) {
	body, err := ioutil.ReadFile(path.Join(metaPath, "tags.txt"))
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return []string{}, err
	}

	return strings.Split(string(body), "\n"), nil
}

func (a *App) WriteFileTags(file string, tags []string) error {
	content := []byte(strings.Join(tags, "\n"))
	metaPath, err := a.MetadataDir(file, false)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(path.Join(metaPath, "tags.txt"), content, 0600)

	return err
}

func (a *App) AddFileTags(file string, tags []string) error {
	origTags, err := a.ReadFileTags(file)
	if err != nil {
		return err
	}
	allTags := append(origTags, tags...)

	return a.WriteFileTags(file, allTags)
}
