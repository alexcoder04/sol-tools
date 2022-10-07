package main

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/alexcoder04/friendly"
)

const VERSION = "0.0.1"

func GetUrl() string {
	return fmt.Sprintf("https://github.com/alexcoder04/sol-lib/archive/refs/tags/v%s.zip", VERSION)
}

func GetLibrary() (string, error) {
	url := GetUrl()
	zipFile := path.Join(os.TempDir(), "sol.zip")
	libFolder := strings.TrimRight(zipFile, ".zip")

	if friendly.IsDir(libFolder) {
		return path.Join(libFolder, fmt.Sprintf("sol-lib-%s", VERSION)), nil
	}

	fmt.Printf("Downloading sol-lib v%s...\n", VERSION)
	err := friendly.DownloadFile(url, zipFile)
	if err != nil {
		return "", err
	}

	fmt.Printf("Extracting sol-lib v%s...\n", VERSION)
	err = friendly.UncompressFolder(zipFile, libFolder)
	if err != nil {
		return "", err
	}

	err = os.Remove(zipFile)
	if err != nil {
		return "", err
	}

	return path.Join(libFolder, fmt.Sprintf("sol-lib-%s", VERSION)), nil
}
