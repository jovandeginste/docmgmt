package app

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"path"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/pkg/browser"
)

type Info struct {
	Filename string `gorm:"type:varchar(512);unique_index;not null"`
	Body     *Body
	Metadata *Metadata
	Info     *FileMetadata
	Tags     Tags `gorm:"type:byte[]"`
	gorm.Model
	App *App `gorm:"-"`
}

type Body struct {
	InfoID  uint
	Content string
	gorm.Model
}

type Metadata struct {
	InfoID uint
	Sha256 string
	Tika   TikaMetadata `gorm:"type:byte[]"`
	gorm.Model
}

type FileMetadata struct {
	InfoID uint
	MTime  int64
	CTime  int64
	Size   int64
	gorm.Model
}

type TikaMetadata map[string][]string
type Tags []string

func (i *Info) IsNew() bool {
	return i.CreatedAt.IsZero()
}

func (i *Info) JSON() []byte {
	infoJSON, _ := json.Marshal(i)

	return infoJSON
}

func (m *Metadata) JSON() []byte {
	metaJSON, _ := json.Marshal(m)

	return metaJSON
}

func (i *Info) AddTag(tag string) {
	for _, t := range i.Tags {
		if t == tag {
			return
		}
	}

	i.Tags = append(i.Tags, tag)
}

func (i *Info) AbsoluteFilename() string {
	return path.Join(i.App.Configuration.DocumentRoot, i.Filename)
}

func (i *Info) Suggestions() ClassificationList {
	suggestions := i.App.Classify(i.Body.Content)

	var remainingSuggestions ClassificationList

outer:
	for _, s := range suggestions {
		for _, t := range i.Tags {
			if t == string(s.Class) {
				continue outer
			}
		}
		remainingSuggestions = append(remainingSuggestions, s)
	}

	return remainingSuggestions
}

func (i *Info) OpenWithDefaultApp() error {
	return browser.OpenFile(i.AbsoluteFilename())
}

func (i *Info) DeleteTag(tag string) {
	var filteredTags []string

	for _, t := range i.Tags {
		if t != tag {
			filteredTags = append(filteredTags, t)
		}
	}

	i.Tags = filteredTags
}

func (i *Info) AddTags(tags []string) {
	for _, t := range tags {
		i.AddTag(t)
	}
}

func (i *Info) Write() error {
	err := i.App.DB.Save(i).Error
	return err
}

func (t *Tags) Scan(value interface{}) error {
	b, ok := value.(string)
	if !ok {
		return fmt.Errorf("Invalid Value")
	}
	if strings.TrimSpace(b) == "" {
		return nil
	}
	*t = strings.Split(strings.TrimSpace(b), "\n")
	return nil
}

func (t Tags) Value() (driver.Value, error) {
	return strings.Join(t, "\n"), nil
}

func (t *TikaMetadata) Scan(value interface{}) error {
	b, ok := value.(string)
	if !ok {
		return fmt.Errorf("Invalid Value")
	}

	err := json.Unmarshal([]byte(b), t)
	return err
}

func (t TikaMetadata) Value() (driver.Value, error) {
	jsonRes, _ := json.Marshal(&t)
	return string(jsonRes), nil
}
