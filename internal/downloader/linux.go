//go:build linux && amd64
// +build linux,amd64

package downloader

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/cheggaaa/pb/v3"
)

const (
	jreVersion     = "21.0.7+6"
	baseURL        = "https://github.com/adoptium/temurin21-binaries/releases/download/jdk-21.0.7+6"
	linuxJRESHA256 = "6d48379e00d47e6fdd417e96421e973898ac90765ea8ff2d09ae0af6d5d6a1c6"

	dynamoURL    = "https://d1ni2b6xgvw0s0.cloudfront.net/v2.x/dynamodb_local_latest.tar.gz"
	dynamoSHA256 = "9a8e6c1b1d4f5c1030c00a5a7eaee1a9ab2b8f1bbde7b700d5505898a3948fff"

	out_jdk    = "jre.tar.gz"
	out_dynamo = "dynamo.tar.gz"
)

var linuxJREURL = baseURL + "/OpenJDK21U-jre_x64_linux_hotspot_21.0.7_6.tar.gz"

func getJREURL() (string, string, string) {
	return linuxJREURL, linuxJRESHA256, out_jdk
}

func getDynamoURL() (string, string, string) {
	return dynamoURL, dynamoSHA256, out_dynamo
}

func Decompress(path, dist string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return err
	}

	bar := pb.Full.Start64(stat.Size())
	barReader := bar.NewProxyReader(file)

	gzReader, err := gzip.NewReader(barReader)
	if err != nil {
		return err
	}
	defer gzReader.Close()

	tr := tar.NewReader(gzReader)
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		target := filepath.Join(dist, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, os.FileMode(header.Mode)); err != nil {
				return err
			}
		case tar.TypeReg:
			if err := os.MkdirAll(filepath.Dir(target), header.FileInfo().Mode().Perm()); err != nil {
				return err
			}
			outFile, err := os.Create(target)
			if err != nil {
				return err
			}

			if _, err := io.Copy(outFile, tr); err != nil {
				outFile.Close()
				return err
			}
			outFile.Close()

			if err := os.Chmod(target, os.FileMode(header.Mode)); err != nil {
				return err
			}
		case tar.TypeSymlink:
			// fmt.Printf("Skipping symlink: %s -> %s\n", header.Name, header.Linkname)
		case tar.TypeLink:
			// fmt.Printf("Skipping hard link: %s -> %s\n", header.Name, header.Linkname)
		case tar.TypeXHeader, tar.TypeGNULongName, tar.TypeGNULongLink, tar.TypeXGlobalHeader:
		default:
			fmt.Printf("Unknown type %v — treating %s as file\n", header.Typeflag, header.Name)
			if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
				return err
			}
			outFile, err := os.Create(target)
			if err != nil {
				return err
			}
			if _, err := io.Copy(outFile, tr); err != nil {
				outFile.Close()
				return err
			}
			outFile.Close()
			if err := os.Chmod(target, os.FileMode(header.Mode)); err != nil {
				return err
			}
		}
	}

	bar.Finish()

	return nil
}
