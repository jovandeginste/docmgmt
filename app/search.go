package app

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"

	"github.com/thedevsaddam/gojsonq"
)

type Result struct {
	Element *Info
	Queries []Query
}

type Query struct {
	Key      string
	Operator string
	Value    string
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
	splitQueries := []Query{}

	for _, q := range queries {
		var a []string
		err := json.Unmarshal([]byte(q), &a)
		if err != nil {
			return []Result{}, err
		}
		theQ := Query{
			Key:      a[0],
			Operator: a[1],
			Value:    a[2],
		}
		splitQueries = append(splitQueries, theQ)

	}

	var result []Result

	for _, f := range files {
		i, err := a.ReadFileInfoMetaPath(f)
		if err != nil {
			continue
		}
		i.Body = ""

		fullJSON := "[" + string(i.JSON()) + "]"

		res := gojsonq.New().JSONString(string(fullJSON))

		for _, q := range splitQueries {
			res = res.Where(q.Key, q.Operator, q.Value)
		}
		theRes := res.Count()

		if theRes > 0 {
			result = append(result, Result{Element: i})
		}
	}

	return result, nil
}
