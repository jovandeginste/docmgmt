package app

import (
	"os"
	"sort"

	"github.com/jovandeginste/docmgmt/bayes"
	"github.com/mitchellh/go-homedir"
)

type Classification struct {
	Class bayes.Class
	Score float64
}

type ClassificationList []Classification

func (p ClassificationList) Len() int           { return len(p) }
func (p ClassificationList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p ClassificationList) Less(i, j int) bool { return p[i].Score > p[j].Score }

func (a *App) Classify(body string) (result ClassificationList) {
	if a.Classifier == nil {
		return
	}

	content := bayes.PrepareString(body)

	x, likely, _ := a.Classifier.LogScores(content)
	a.Logf(LogDebug, "Classification top hit: '%s'; score: '%f'", a.Classifier.Classes[likely], x[likely])

	total := float64(0)

	for i, c := range a.Classifier.Classes {
		score := -1 / x[i]

		result = append(result, Classification{Class: c, Score: score})

		total += score
	}

	for i := range result {
		result[i].Score /= total
	}

	sort.Sort(result)

	return
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
