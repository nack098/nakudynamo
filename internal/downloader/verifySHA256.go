package downloader

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
)

func verifySHA256(filePath, expected string) (bool, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return false, err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return false, err
	}

	actual := hex.EncodeToString(h.Sum(nil))
	return actual == expected, nil
}
