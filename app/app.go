package app

import (
	"context"
	"errors"
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
	a.Logf(LogInfo, "Starting tika server version: "+a.Configuration.TikaVersion)

	return a.TikaServer.Start(context.Background())
}

func (a *App) StopServer() {
	a.Logf(LogDebug, "Stopping tika server")

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
	return a.ReadAbsoluteFileInfo(path)
}

func (a *App) ReadAbsoluteFileInfo(file string) (result *Info, err error) {
	result = &Info{Filename: file}
	result.App = a

	err = a.DB.Where(result).FirstOrInit(&result).Error
	return
}

func (a *App) AllTags() []string {
	var (
		result []string
		tags   string
	)

	f := map[string]bool{}

	rows, _ := a.DB.Table("infos").Select("tags").Rows() //nolint:errcheck
	defer rows.Close()

	for rows.Next() {
		rows.Scan(&tags) //nolint:errcheck

		for _, t := range strings.Split(tags, "\n") {
			if _, ok := f[t]; !ok {
				f[t] = true

				result = append(result, t)
			}
		}
	}

	return result
}
