package app

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/elgs/jsonql"
)

type Result struct {
	Element *Info
}

func (a *App) MetadataFiles() []string {
	result := []string{}

	err := filepath.Walk(a.Configuration.MetadataRoot,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if filepath.Base(path) == "metadata.json" {
				result = append(result, filepath.Dir(path))
			}
			return nil
		})
	if err != nil {
		log.Println(err)
	}

	return result
}

func (a *App) Search(queries []string) ([]Result, error) {
	files := a.MetadataFiles()
	fullQuery := strings.Join(queries, " ")

	var result []Result

	for _, f := range files {
		i, err := a.ReadFileInfoMetaPath(f)
		if err != nil {
			continue
		}

		fullJSON := string(i.JSON())
		parser, err := jsonql.NewStringQuery(fullJSON)
		if err != nil {
			return result, err
		}
		res, err := parser.Query(fullQuery)
		if err != nil {
			return result, err
		}
		jsonRes, _ := json.Marshal(res)

		if string(jsonRes) != "null" {
			result = append(result, Result{Element: i})
		}
	}

	return result, nil
}
