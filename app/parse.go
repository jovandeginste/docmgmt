package app

import (
	"bytes"
	"context"
	"encoding/csv"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/google/go-tika/tika"
)

func (a *App) fileMetadata(file string) (*FileMetadata, error) {
	a.Logf(LogDebug, "Parsing filesystem metadata for: '%s'", file)

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

func (a *App) Parse(file string) (info *Info, err error) {
	a.Logf(LogInfo, "Parsing file: '%s'", file)

	info, err = a.ReadFileInfo(file)
	if err != nil {
		return
	}

	fileMeta, err := a.fileMetadata(file)
	if err != nil {
		return
	}

	info.Info = fileMeta

	a.Logf(LogDebug, "Reading content of: '%s'", file)

	f, err := os.Open(file)
	if err != nil {
		return
	}
	defer f.Close()

	content, err := ioutil.ReadAll(f)
	if err != nil {
		return
	}

	info.Metadata = &Metadata{
		Sha256: checksum(content),
		Tika:   TikaMetadata{},
	}
	info.Body = &Body{
		Content: string(content),
	}

	client := tika.NewClient(nil, a.TikaServer.URL())

	err = info.ParseMetadata(client)
	if err != nil {
		return
	}

	err = info.ParseBody(client)

	return info, err
}

func (i *Info) ParseMetadata(client *tika.Client) error {
	i.App.Logf(LogDebug, "Fetching Tika metadata for: '%s'", i.Filename)

	meta, err := client.Meta(context.Background(), strings.NewReader(i.Body.Content))
	if err != nil {
		return err
	}

	for _, s := range strings.Split(strings.TrimSpace(meta), "\n") {
		r := csv.NewReader(strings.NewReader(s))
		fields, _ := r.Read()

		i.Metadata.Tika[fields[0]] = fields[1:]
	}

	return nil
}

func (i *Info) ParseBody(client *tika.Client) error {
	i.App.Logf(LogDebug, "Fetching Tika body for: '%s'", i.Filename)

	body, err := client.Parse(context.Background(), strings.NewReader(i.Body.Content))
	if err != nil {
		return err
	}

	if body = strings.TrimSpace(body); body != "" {
		i.Body.Content = body

		return nil
	}

	if len(i.Metadata.Tika["Content-Type"]) == 1 && i.Metadata.Tika["Content-Type"][0] == "application/pdf" {
		err = i.ParseBodyAsImage(client)
		if err != nil {
			return err
		}
	}

	return err
}

func (i *Info) ParseBodyAsImage(client *tika.Client) (err error) {
	nPages := 1

	if len(i.Metadata.Tika["xmpTPg:NPages"]) > 0 {
		nPagesStr := i.Metadata.Tika["xmpTPg:NPages"][0]
		nPages, err = strconv.Atoi(nPagesStr)

		if err != nil {
			return
		}
	}

	i.App.Logf(LogInfo, "PDF body seems empty; trying to convert OCR. We seem to have %d page(s).", nPages)

	var (
		results     []image
		body, bodyN string
	)

	results, err = convertPdfToImage(image(i.Body.Content), nPages)
	if err != nil {
		return
	}

	i.App.Logf(LogInfo, "Performing OCR scan on images.")

	for _, r := range results {
		bodyN, err = client.Parse(context.Background(), bytes.NewReader(r))
		if err != nil {
			return
		}

		body += bodyN + "\n"
	}

	i.Body.Content = strings.TrimSpace(body)

	return
}
