package app

import (
	"context"
	"errors"
	"log"
	"path/filepath"
	"strings"

	"github.com/google/go-tika/tika"
	"github.com/jinzhu/gorm"
	"github.com/jovandeginste/docmgmt/bayes"
)

type App struct {
	Configuration Configuration
	DB            *gorm.DB
	Classifier    *bayes.Classifier
	TikaServer    *tika.Server
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
	path, err := filepath.Abs(file)
	if err != nil {
		return
	}

	if !strings.HasPrefix(path, a.Configuration.DocumentRoot) {
		err = errors.New("file is not inside document root")
		return
	}

	path = strings.TrimPrefix(path, a.Configuration.DocumentRoot)

	result = &Info{Filename: path}
	result.App = a

	err = a.DB.Where(result).FirstOrInit(&result).Error
	return
}
