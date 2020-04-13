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
	var (
		result   []Result
		filename string
	)

	fullQuery := strings.Join(queries, " ")

	rows, _ := a.DB.Table("infos").Select("filename").Rows() //nolint:errcheck
	defer rows.Close()

	for rows.Next() {
		rows.Scan(&filename) //nolint:errcheck

		info, err := a.ReadAbsoluteFileInfo(filename)
		if err != nil {
			continue
		}

		fullJSON := string(info.JSON())
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
			result = append(result, Result{Element: info})
		}
	}

	return result, nil
}
