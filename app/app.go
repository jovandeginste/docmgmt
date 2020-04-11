package app

import (
	"context"
	"log"

	"github.com/google/go-tika/tika"
	"github.com/jinzhu/gorm"
	"github.com/jovandeginste/docmgmt/bayes"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

type App struct {
	Configuration Configuration
	DB            *gorm.DB
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

	a.TikaServer = s

	return nil
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

func (a *App) ReadFileInfo(file string) (result *Info, err error) {
	result = &Info{Filename: file}
	result.App = a

	err = a.DB.Where(result).FirstOrInit(&result).Error
	return
}
