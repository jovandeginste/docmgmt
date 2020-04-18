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

	if strings.TrimSpace(body) == "" && len(metaResult.Tika["Content-Type"]) == 1 && metaResult.Tika["Content-Type"][0] == "application/pdf" {
		nPages := 1
		if len(metaResult.Tika["xmpTPg:NPages"]) > 0 {
			nPagesStr := metaResult.Tika["xmpTPg:NPages"][0]
			nPages, err = strconv.Atoi(nPagesStr)
			if err != nil {
				return
			}
		}

		a.Logf(LogInfo, "PDF body seems empty; trying to convert OCR. We seem to have %d page(s).", nPages)

		var results [][]byte

		results, err = ConvertPdfToJpg(file, nPages)
		if err != nil {
			return
		}

		a.Logf(LogInfo, "Performing OCR scan on images.")
		body = ""
		bodyN := ""
		for _, r := range results {
			bodyN, err = client.Parse(context.Background(), bytes.NewReader(r))
			body += bodyN + "\n"
		}

		info.Body = &Body{Content: body}
	}

	return
}
