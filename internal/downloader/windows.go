//go:build windows && amd64
// +build windows,amd64

package downloader

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"

	"github.com/cheggaaa/pb/v3"
)

const (
	jreVersion       = "21.0.7+6"
	baseURL          = "https://github.com/adoptium/temurin21-binaries/releases/download/jdk-21.0.7+6"
	windowsJRESHA256 = "b2850a96293048ed3020f8bfca2d92a785ae9bf80c7d96bbfe3ec4ccf45aef98"

	dynamoURL    = "https://d1ni2b6xgvw0s0.cloudfront.net/v2.x/dynamodb_local_latest.zip"
	dynamoSHA256 = "06e7bdd5d03262d8373696282f79867f9f6beb94b76e230b49da135be080558c"

	out_jdk    = "jre.zip"
	out_dynamo = "dynamo.zip"
)

var windowsJREURL = baseURL + "/OpenJDK21U-jre_x64_windows_hotspot_21.0.7_6.zip"

func getJREURL() (string, string, string) {
	return windowsJREURL, windowsJRESHA256, out_jdk
}

func getDynamoURL() (string, string, string) {
	return dynamoURL, dynamoSHA256, out_dynamo
}

func Decompress(path, dist string) error {
	r, err := zip.OpenReader(path)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		target := filepath.Join(dist, f.Name)

		if f.FileInfo().IsDir() {
			err := os.MkdirAll(target, f.Mode())
			if err != nil {
				return err
			}
			continue
		}

		if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
			return err
		}

		rc, err := f.Open()
		if err != nil {
			return err
		}

		outFile, err := os.OpenFile(target, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			rc.Close()
			return err
		}

		bar := pb.Full.Start64(int64(f.UncompressedSize64))
		barReader := bar.NewProxyReader(rc)

		_, err = io.Copy(outFile, barReader)
		if err != nil {
			outFile.Close()
			rc.Close()
			return err
		}

		bar.Finish()

		outFile.Close()
		rc.Close()
	}

	return nil
}
