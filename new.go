package main

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/alexcoder04/friendly"
)

func New(projectName string) error {
	for _, c := range []string{"/", ".", "\\"} {
		if strings.Contains(projectName, c) {
			return errors.New("invalid project name")
		}
	}

	if friendly.Exists(projectName) {
		return os.ErrExist
	}

	folders := []string{
		projectName,
		path.Join(projectName, "components"),
		path.Join(projectName, "res"),
		path.Join(projectName, "data")}
	for _, f := range folders {
		err := os.Mkdir(f, 0700)
		if err != nil {
			return err
		}
	}

	err := friendly.WriteNewFile(path.Join(projectName, "res", "data", "menu.yml"), "[]")
	if err != nil {
		return err
	}

	exe, err := os.Executable()
	if err != nil {
		return err
	}

	err = friendly.WriteNewFile(path.Join(projectName, "Makefile"), fmt.Sprintf(`
SOL = %s/sol
UPLOADNSPIRE = uploadnspire
NAME = helloworld
TEMP_LUA = %s/out.lua
OUT_FILE = %s/$(NAME).tns

all: clean build upload

build:
	$(SOL) .
	luna $(TEMP_LUA) $(OUT_FILE)

clean:
	$(RM) $(TEMP_LUA) $(OUT_FILE)

upload:
	$(UPLOADNSPIRE) $(OUT_FILE)
`, exe, os.TempDir(), os.TempDir()))
	if err != nil {
		return err
	}

	err = friendly.WriteNewFile(path.Join(projectName, "README.md"), "# helloworld application for the ti-nspire")
	if err != nil {
		return err
	}

	err = friendly.WriteNewFile(path.Join(projectName, "app.lua"), `
hello_world_element = Components.Base.TextField:new()
hello_world_element.Label = "Hello World"

App:AddElement(hello_world_element)
`)
	if err != nil {
		return err
	}

	return nil
}
