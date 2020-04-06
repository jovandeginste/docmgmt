package app

import (
	"bytes"
	"context"
	"encoding/csv"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"syscall"
	"time"

	"github.com/google/go-tika/tika"
)

type Metadata struct {
	Filename string
	Info     *FileMetadata
	Tika     TikaMetadata
}

type FileMetadata struct {
	MTime int64
	CTime int64
	Size  int64
}

type TikaMetadata map[string][]string

func fileMetadata(file string) (*FileMetadata, error) {
	s, err := os.Stat(file)
	if err != nil {
		return nil, err
	}

	mtime := s.ModTime()
	stat := s.Sys().(*syscall.Stat_t)
	ctime := time.Unix(stat.Ctim.Sec, stat.Ctim.Nsec)

	result := FileMetadata{
		MTime: mtime.UnixNano(),
		CTime: ctime.UnixNano(),
		Size:  s.Size(),
	}

	return &result, nil
}

func (a *App) Parse(file string) (*Metadata, string, error) {
	if a.Configuration.Verbose {
		log.Printf("Parsing file: %#v", file)
	}

	fileMeta, err := fileMetadata(file)
	if err != nil {
		return nil, "", err
	}

	f, err := os.Open(file)
	if err != nil {
		return nil, "", err
	}
	defer f.Close()

	content, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, "", err
	}

	client := tika.NewClient(nil, a.TikaServer.URL())

	meta, err := client.Meta(context.Background(), bytes.NewReader(content))
	if err != nil {
		return nil, "", err
	}

	metaResult := Metadata{
		Filename: file,
		Info:     fileMeta,
		Tika:     TikaMetadata{},
	}

	for _, s := range strings.Split(strings.TrimSpace(meta), "\n") {
		r := csv.NewReader(strings.NewReader(s))
		fields, _ := r.Read()

		metaResult.Tika[fields[0]] = fields[1:]
	}

	body, err := client.Parse(context.Background(), bytes.NewReader(content))

	return &metaResult, body, err
}
