package app

import (
	"fmt"

	"github.com/jovandeginste/docmgmt/bayes"
)

func (a *App) Classify(body string) string {
	if a.Classifier == nil {
		return ""
	}

	content := bayes.PrepareString(body)

	_, likely, _ := a.Classifier.LogScores(content)

	fmt.Println(likely, a.Classifier.Classes[likely])

	return string(a.Classifier.Classes[likely])
}
