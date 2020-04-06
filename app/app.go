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

	"github.com/google/go-tika/tika"
	"github.com/jovandeginste/docmgmt/bayes"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
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

	if a.Configuration.Verbose {
		log.Println("Tika server version: " + a.Configuration.TikaVersion)
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
		log.Println("Starting tika server")
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
	metaPath := a.MetadataDir(file)

	if a.Configuration.Verbose {
		log.Println("Metadata dir: " + metaPath)
	}

	meta, body, err := a.Parse(file)
	if err != nil {
		return err
	}

	metaJSON, err := json.Marshal(meta)
	os.MkdirAll(metaPath, 0700)

	ioutil.WriteFile(path.Join(metaPath, "metadata.json"), metaJSON, 0600)
	ioutil.WriteFile(path.Join(metaPath, "body.txt"), []byte(body), 0600)
	return nil
}

func (a *App) MetadataDir(file string) string {
	sum := sha256.Sum256([]byte(file))
	relPath := fmt.Sprintf("%x", sum)

	return path.Join(a.Configuration.MetadataRoot, relPath[0:2], relPath[2:4], relPath)
}

func (a *App) ReadFileBody(file string) (string, error) {
	metaPath := a.MetadataDir(file)
	body, err := ioutil.ReadFile(path.Join(metaPath, "body.txt"))
	if !os.IsNotExist(err) {
		return "", errors.New("file not parsed yet")
	}

	return string(body), err
}

func (a *App) ReadFileMetadata(file string) (string, error) {
	metaPath := a.MetadataDir(file)
	meta, err := ioutil.ReadFile(path.Join(metaPath, "metadata.json"))
	if !os.IsNotExist(err) {
		return "", errors.New("file not parsed yet")
	}

	return string(meta), err
}
