package app

import (
	"encoding/json"
	"strings"

	"github.com/elgs/jsonql"
)

type Result struct {
	Element *Info
}

func (a *App) Search(queries []string) ([]Result, error) {
	var files []string
	a.DB.Table("infos").Pluck("filename", &files)

	fullQuery := strings.Join(queries, " ")

	var result []Result

	for _, f := range files {
		i, err := a.ReadAbsoluteFileInfo(f)
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
