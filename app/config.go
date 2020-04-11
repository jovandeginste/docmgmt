package app

import (
	"strings"

	"github.com/google/go-tika/tika"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

type Configuration struct {
	Verbose        bool
	DocumentRoot   string `mapstructure:"document_root"`
	MetadataDB     string `mapstructure:"metadata_db"`
	TikaVersion    string `mapstructure:"tika_version"`
	ClassifierData string `mapstructure:"classifier_data"`
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

	err = a.InitDB()
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

	if !strings.HasSuffix(a.Configuration.DocumentRoot, "/") {
		a.Configuration.DocumentRoot = a.Configuration.DocumentRoot + "/"
	}

	a.TikaServer = s

	return nil
}
