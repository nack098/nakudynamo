package downloader

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/cheggaaa/pb/v3"
)

func DownloadDynamo(destDir string) (string, error) {
	url, hash, out := getDynamoURL()

	outPath := filepath.Join(destDir, out)
	if _, err := os.Stat(outPath); err == nil {
		ok, err := verifySHA256(outPath, hash)
		if err != nil {
			return "", fmt.Errorf("failed to verify existing DynamoDBLocal file: %w", err)
		}
		if ok {
			fmt.Println("DynamoDBLocal archive already downloaded and verified.")
			return outPath, nil
		}
		fmt.Println("Checksum mismatch, re-downloading DynamoDBLocal...")
	}

	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to download DynamoDBLocal: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("bad response downloading DynamoDBLocal: %s", resp.Status)
	}

	outFile, err := os.Create(outPath)
	if err != nil {
		return "", fmt.Errorf("failed to create file %s: %w", outPath, err)
	}
	defer outFile.Close()

	bar := pb.Full.Start64(resp.ContentLength)
	barReader := bar.NewProxyReader(resp.Body)

	if _, err := io.Copy(outFile, barReader); err != nil {
		return "", fmt.Errorf("failed to write DynamoDBLocal file: %w", err)
	}
	bar.Finish()

	ok, err := verifySHA256(outPath, hash)
	if err != nil {
		return "", fmt.Errorf("failed to verify DynamoDBLocal checksum after download: %w", err)
	}
	if !ok {
		os.Remove(outPath)
		return "", fmt.Errorf("checksum mismatch after download")
	}

	fmt.Println("✅ DynamoDBLocal downloaded and verified successfully!")
	return outPath, nil
}
