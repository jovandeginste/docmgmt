package app

import (
	"fmt"
	"os"

	"github.com/jovandeginste/docmgmt/bayes"
	"github.com/mitchellh/go-homedir"
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

func (a *App) Learn(body string, tag string) {
	tagClass := bayes.Class(tag)
	content := bayes.PrepareString(body)

	a.Classifier.AddClass(tagClass)
	a.Classifier.Learn(content, tagClass)
}

func (a *App) LoadClassifier() error {
	cp, err := homedir.Expand(a.Configuration.ClassifierData)
	if err != nil {
		return err
	}

	c, err := bayes.NewClassifierFromFile(cp)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
		c = bayes.NewClassifier()
	}

	a.Classifier = c
	return nil
}

func (a *App) SaveClassifier() error {
	cp, err := homedir.Expand(a.Configuration.ClassifierData)
	if err != nil {
		return err
	}

	return a.Classifier.WriteToFile(cp)
}
