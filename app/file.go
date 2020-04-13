package app

import (
	"crypto/sha256"
	"fmt"
)

func checksum(content []byte) (checksumSha256 string) {
	sha256Sum := sha256.Sum256(content)
	checksumSha256 = fmt.Sprintf("%x", sha256Sum)

	return
}
