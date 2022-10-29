package commands

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/alexcoder04/arrowprint"
	"github.com/alexcoder04/friendly/v2/ffiles"
	"github.com/alexcoder04/sol-tools/utils"
)

func New(pfolder string, projectName string) error {
	arrowprint.InfoC("Creating new project")
	arrowprint.Info0("Doing initial checks")
	arrowprint.Info1("Checking project name")
	for _, c := range []string{"/", ".", "\\"} {
		if strings.Contains(projectName, c) {
			return errors.New("invalid project name")
		}
	}

	folder := path.Join(pfolder, projectName)

	arrowprint.Info1("Checking if project already exists")
	exists, err := ffiles.Exists(folder)
	if err != nil {
		return err
	}
	if exists {
		return os.ErrExist
	}

	arrowprint.Info0("Creating folder structure")
	folders := []string{
		folder,
		path.Join(folder, "components"),
		path.Join(folder, "res"),
		path.Join(folder, "res/data")}
	for _, f := range folders {
		arrowprint.Info1("Creating %s", f)
		err := os.Mkdir(f, 0700)
		if err != nil {
			return err
		}
	}

	arrowprint.Info0("Generating boilerplate code")
	arrowprint.Info1("Creating menu.yml file")
	err = ffiles.WriteNewFile(path.Join(folder, "res", "data", "menu.yml"), "[]")
	if err != nil {
		return err
	}

	arrowprint.Warn1("Generating Makefile")
	err = ffiles.WriteNewFile(path.Join(folder, "Makefile"), fmt.Sprintf(`
NAME = %s
TEMP_LUA = %s/out.lua
OUT_FILE = %s/$(NAME).tns

all: clean build upload

build:
	sol -a build .
	luna $(TEMP_LUA) $(OUT_FILE)

clean:
	$(RM) $(TEMP_LUA) $(OUT_FILE)

upload:
	uploadnspire $(OUT_FILE)
`, projectName, os.TempDir(), os.TempDir()))
	if err != nil {
		return err
	}

	arrowprint.Info1("Generating README.md file")
	err = ffiles.WriteNewFile(
		path.Join(folder, "README.md"),
		fmt.Sprintf("# %s\nAn app for the ti-nspire", projectName))
	if err != nil {
		return err
	}

	arrowprint.Info1("Creating app.lua file")
	err = ffiles.WriteNewFile(path.Join(folder, "app.lua"), `
hello_world_element = Components.Base.TextField:new()
hello_world_element.Label = "Hello World"

App:AddElement(hello_world_element)
`)
	if err != nil {
		return err
	}

	arrowprint.Info1("Creating solproj.yml file")
	version, err := utils.GetLatestVersion()
	if err != nil {
		return err
	}
	err = ffiles.WriteNewFile(path.Join(folder, "solproj.yml"), fmt.Sprintf(`
RefreshRate: 0.5
SolVersion: %s
`, version))
	if err != nil {
		return err
	}

	arrowprint.Info1("Creating .gitignore file")
	err = ffiles.WriteNewFile(path.Join(folder, ".gitignore"), `
out.lua
*.tns
`)
	if err != nil {
		return err
	}

	arrowprint.Suc0("Your project '%s' has been set up successfully", projectName)
	return nil
}
