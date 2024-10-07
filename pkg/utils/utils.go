package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/thiagozs/go-sftpsync/internal/domain"
	"gopkg.in/yaml.v3"
)

func GetFileInfo(file string) (os.FileInfo, error) {
	fopen, err := os.Open(file)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}

	fileInfo, err := fopen.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to get file size: %v", err)
	}

	return fileInfo, nil
}

func GetCurrentDirectory() (string, error) {
	currentPath, err := filepath.Abs(filepath.Dir("."))
	if err != nil {
		return "", err
	}
	return currentPath, nil
}

func ReadConfig(filePath string) (*domain.Config, error) {
	// Check if the file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file not found: %s", filePath)
	}

	// Read the file content
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	// Unmarshal the YAML into the Config struct
	var config domain.Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling YAML: %w", err)
	}

	return &config, nil
}

func ParsePath(fullPath, syncFolder string) (string, string, error) {
	// Ensure the full path is clean (removes redundant slashes, etc.)
	fullPath = filepath.Clean(fullPath)

	// Find the position of the syncFolder in the full path
	index := strings.Index(fullPath, syncFolder)
	if index == -1 {
		return "", "", fmt.Errorf("sync folder %s not found in path", syncFolder)
	}

	// Split the path into two chunks
	// First part: everything up to the syncFolder
	// Second part: everything after the syncFolder
	firstPart := fullPath[:index+len(syncFolder)]
	secondPart := fullPath[index+len(syncFolder):]

	return firstPart, secondPart, nil
}
