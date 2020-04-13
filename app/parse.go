package app

import (
	"bytes"
	"context"
	"encoding/csv"
	"io/ioutil"
	"os"
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

	metaResult := Metadata{
		Sha256: checksum(content),
		Tika:   TikaMetadata{},
	}
	info.Metadata = &metaResult

	a.Logf(LogDebug, "Fetching Tika metadata for: '%s'", file)
	client := tika.NewClient(nil, a.TikaServer.URL())
	meta, err := client.Meta(context.Background(), bytes.NewReader(content))
	if err != nil {
		return
	}

	for _, s := range strings.Split(strings.TrimSpace(meta), "\n") {
		r := csv.NewReader(strings.NewReader(s))
		fields, _ := r.Read()

		metaResult.Tika[fields[0]] = fields[1:]
	}

	a.Logf(LogDebug, "Fetching Tika body for: '%s'", file)
	body, err := client.Parse(context.Background(), bytes.NewReader(content))
	info.Body = &Body{Content: body}

	return
}
