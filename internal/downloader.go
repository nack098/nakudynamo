package internal

import (
	"fmt"
	"nakuya/nakudynamo/internal/downloader"
	"os"
	"path/filepath"
	"runtime"
)

func renameExtractedFolder(dist, oldName, newName string) error {
	oldPath := filepath.Join(dist, oldName)
	newPath := filepath.Join(dist, newName)
	return os.Rename(oldPath, newPath)
}

func PrepareEnvironment() (*DynamoEnvironment, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("could not get user home: %w", err)
	}

	workingDir := filepath.Join(home, ".nakudynamo")
	download_dir := filepath.Join(workingDir, ".tmp")
	jrePath := filepath.Join(workingDir, "jre", "bin", "java")
	if runtime.GOOS == "windows" {
		jrePath += ".exe"
	}
	jarPath := filepath.Join(workingDir, "DynamoDBLocal.jar")

	if err := os.MkdirAll(workingDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create working dir: %w", err)
	}

	if err := os.MkdirAll(download_dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create working dir: %w", err)
	}

	if _, err := os.Stat(jrePath); os.IsNotExist(err) {
		jreDownloadPath, err := downloader.DownloadJRE(download_dir)
		if err != nil {
			return nil, fmt.Errorf("failed to download jre: %w", err)
		}
		downloader.Decompress(jreDownloadPath, workingDir)
		if err := renameExtractedFolder(workingDir, "jdk-21.0.7+6-jre", "jre"); err != nil {
			return nil, fmt.Errorf("cannot rename the folder: %w", err)
		}
	}

	if _, err := os.Stat(jarPath); os.IsNotExist(err) {
		dynamoDownloadPath, err := downloader.DownloadDynamo(download_dir)
		if err != nil {
			return nil, fmt.Errorf("failed to download jre: %w", err)
		}

		downloader.Decompress(dynamoDownloadPath, workingDir)
	}

	return &DynamoEnvironment{
		JREPath:       jrePath,
		DynamoJarPath: jarPath,
		WorkingDir:    workingDir,
		Port:          8000,
	}, nil
}
