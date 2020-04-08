package app

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/thedevsaddam/gojsonq"
)

type Result struct {
	Element *Info
	Queries []Query
}
type Query struct {
	Q string
	V string
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
	result := []Result{}
	splitQueries := []Query{}

	for _, q := range queries {
		splitQ := strings.SplitN(q, "=", 2)
		if len(splitQ) > 1 {
			splitQueries = append(splitQueries, Query{splitQ[0], splitQ[1]})
		}
	}

	for _, f := range files {
		i, err := a.ReadFileInfoMetaPath(f)
		if err != nil {
			continue
		}

		e := Result{
			Element: i,
			Queries: []Query{},
		}

		match := true
		res := gojsonq.New().JSONString(string(i.JSON()))

		for _, q := range splitQueries {
			qr, err := res.FindR(q.Q)

			if err != nil {
				match = false
				continue
			}

			names, err := qr.StringSlice()

			if err != nil {
				name, err := qr.String()
				if err != nil {
					match = false
					continue
				}

				names = []string{name}
			}

			subMatch := false
			for _, name := range names {
				if name == q.V {
					subMatch = true
				}
			}

			if !subMatch {
				match = false
			}
		}

		if match {
			result = append(result, e)
		}
	}

	return result, nil
}
