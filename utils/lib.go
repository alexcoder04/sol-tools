package utils

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/alexcoder04/arrowprint"
	"github.com/alexcoder04/friendly"
)

const VERSION = "0.0.1"

func GetUrl() string {
	return fmt.Sprintf("https://github.com/alexcoder04/sol-lib/archive/refs/tags/v%s.zip", VERSION)
}

func GetLibrary() (string, error) {
	if os.Getenv("SOL_USE_LOCAL_LIB") == "1" {
		folder := os.Getenv("SOL_LOCAL_LIB_PATH")
		if !friendly.IsDir(folder) {
			return "", os.ErrNotExist
		}
		return folder, nil
	}

	url := GetUrl()
	zipFile := path.Join(os.TempDir(), "sol.zip")
	libFolder := strings.TrimRight(zipFile, ".zip")

	if friendly.IsDir(libFolder) {
		return path.Join(libFolder, fmt.Sprintf("sol-lib-%s", VERSION)), nil
	}

	arrowprint.Info1("Downloading sol-lib v%s...", VERSION)
	err := friendly.DownloadFile(url, zipFile)
	if err != nil {
		return "", err
	}

	arrowprint.Info1("Extracting sol-lib v%s...", VERSION)
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
