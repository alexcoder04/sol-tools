package main

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/alexcoder04/arrowprint"
	"github.com/alexcoder04/friendly"
)

func New(projectName string) error {
	arrowprint.InfoC("Creating new project")
	arrowprint.Info1("Checking project name")
	for _, c := range []string{"/", ".", "\\"} {
		if strings.Contains(projectName, c) {
			return errors.New("invalid project name")
		}
	}

	arrowprint.Info1("Checking if project already exists")
	if friendly.Exists(projectName) {
		return os.ErrExist
	}

	arrowprint.Info0("Creating folder structure")
	folders := []string{
		projectName,
		path.Join(projectName, "components"),
		path.Join(projectName, "res"),
		path.Join(projectName, "res/data")}
	for _, f := range folders {
		arrowprint.Info1("Creating %s", f)
		err := os.Mkdir(f, 0700)
		if err != nil {
			return err
		}
	}

	arrowprint.Info1("Creating menu.yml file")
	err := friendly.WriteNewFile(path.Join(projectName, "res", "data", "menu.yml"), "[]")
	if err != nil {
		return err
	}

	arrowprint.Warn1("Generating Makefile")
	err = friendly.WriteNewFile(path.Join(projectName, "Makefile"), fmt.Sprintf(`
SOL = sol
UPLOADNSPIRE = uploadnspire
NAME = helloworld
TEMP_LUA = %s/out.lua
OUT_FILE = %s/$(NAME).tns

all: clean build upload

build:
	$(SOL) -a build .
	luna $(TEMP_LUA) $(OUT_FILE)

clean:
	$(RM) $(TEMP_LUA) $(OUT_FILE)

upload:
	$(UPLOADNSPIRE) $(OUT_FILE)
`, os.TempDir(), os.TempDir()))
	if err != nil {
		return err
	}

	arrowprint.Info1("Generating README.md file")
	err = friendly.WriteNewFile(
		path.Join(projectName, "README.md"),
		fmt.Sprintf("# %s\nAn app for the ti-nspire", projectName))
	if err != nil {
		return err
	}

	arrowprint.Info1("Creating app.lua file")
	err = friendly.WriteNewFile(path.Join(projectName, "app.lua"), `
hello_world_element = Components.Base.TextField:new()
hello_world_element.Label = "Hello World"

App:AddElement(hello_world_element)
`)
	if err != nil {
		return err
	}

	arrowprint.Info1("Creating solproj.yml file")
	err = friendly.WriteNewFile(path.Join(projectName, "solproj.yml"), `
RefreshRate = 0.5
SolVersion = 0
`)
	if err != nil {
		return err
	}

	arrowprint.Suc0("Your project %s has been set up successfully", projectName)
	return nil
}
