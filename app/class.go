package app

import (
	"fmt"
	"strings"
)

func (a *App) Classify(content string) string {
	if a.Classifier == nil {
		return ""
	}

	_, likely, _ := a.Classifier.LogScores(
		strings.Fields(strings.ToLower(string(content))),
	)

	fmt.Println(likely, a.Classifier.Classes[likely])

	return string(a.Classifier.Classes[likely])
}
