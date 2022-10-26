package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/alexcoder04/arrowprint"
	"github.com/alexcoder04/friendly"
)

func GetUrl(version string) string {
	return fmt.Sprintf("https://github.com/alexcoder04/sol-lib/archive/refs/tags/v%s.zip", version)
}

func GetLibrary(version string) (string, string, error) {
	if os.Getenv("SOL_USE_LOCAL_LIB") == "1" {
		folder := os.Getenv("SOL_LOCAL_LIB_PATH")
		if !friendly.IsDir(folder) {
			return "", "<local>", os.ErrNotExist
		}
		return folder, "", nil
	}

	if version == "latest" {
		resp, err := http.Get("https://api.github.com/repos/alexcoder04/sol-lib/git/refs/tags")
		if err != nil {
			return "", "", err
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", "", err
		}
		var releasesData []GitHubApiRef
		err = json.Unmarshal(body, &releasesData)
		if err != nil {
			return "", "", err
		}
		latestIndex := 0
		latestVersion := strings.TrimPrefix(releasesData[latestIndex].Ref, "refs/tags/v")
		for i, r := range releasesData {
			v := strings.TrimPrefix(r.Ref, "refs/tags/v")
			if friendly.SemVersionGreater(v, latestVersion) {
				latestVersion = v
				latestIndex = i
			}
		}
		version = latestVersion
		arrowprint.Info1("Found latest sol version: %s", version)
	}

	url := GetUrl(version)
	zipFile := path.Join(os.TempDir(), "sol.zip")
	libFolder := path.Join(strings.TrimRight(zipFile, ".zip"), fmt.Sprintf("sol-lib-%s", version))

	if friendly.IsDir(libFolder) {
		return libFolder, version, nil
	}

	arrowprint.Info1("Downloading sol-lib v%s...", version)
	err := friendly.DownloadFile(url, zipFile)
	if err != nil {
		return "", "", err
	}

	arrowprint.Info1("Extracting sol-lib v%s...", version)
	err = friendly.UncompressFolder(zipFile, path.Dir(libFolder))
	if err != nil {
		return "", "", err
	}

	err = os.Remove(zipFile)
	if err != nil {
		return "", "", err
	}

	return libFolder, version, nil
}
