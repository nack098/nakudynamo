package downloader

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/cheggaaa/pb/v3"
)

func DownloadJRE(destDir string) (string, error) {
	url, hash, out := getJREURL()

	outPath := filepath.Join(destDir, out)
	if _, err := os.Stat(outPath); err == nil {
		ok, err := verifySHA256(outPath, hash)
		if err != nil {
			return "", fmt.Errorf("failed to verify existing JRE file: %w", err)
		}
		if ok {
			fmt.Println("JRE archive already downloaded and verified.")
			return outPath, nil
		}
		fmt.Println("Checksum mismatch, re-downloading JRE...")
	}

	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to download JRE: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("bad response downloading JRE: %s", resp.Status)
	}

	outFile, err := os.Create(outPath)
	if err != nil {
		return "", fmt.Errorf("failed to create file %s: %w", outPath, err)
	}
	defer outFile.Close()

	if resp.ContentLength > 0 {
		bar := pb.Full.Start64(resp.ContentLength)
		barReader := bar.NewProxyReader(resp.Body)
		_, err = io.Copy(outFile, barReader)
		bar.Finish()
	} else {
		fmt.Println("Downloading (unknown size)...")
		_, err = io.Copy(outFile, resp.Body)
	}
	if err != nil {
		return "", fmt.Errorf("failed to write JRE file: %w", err)
	}

	ok, err := verifySHA256(outPath, hash)
	if err != nil {
		return "", fmt.Errorf("failed to verify JRE checksum after download: %w", err)
	}
	if !ok {
		os.Remove(outPath)
		return "", fmt.Errorf("checksum mismatch after download")
	}

	fmt.Println("JRE downloaded and verified successfully!")
	return outPath, nil
}
