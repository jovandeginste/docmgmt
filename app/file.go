package app

import (
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"os"
)

func fileChecksum(file string) (checksumSha256 string, err error) {
	f, err := os.Open(file)
	if err != nil {
		return
	}
	defer f.Close()

	content, err := ioutil.ReadAll(f)
	if err != nil {
		return
	}

	sha256Sum := sha256.Sum256(content)
	checksumSha256 = fmt.Sprintf("%x", sha256Sum)

	return
}
