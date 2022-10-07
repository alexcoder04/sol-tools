package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/alexcoder04/friendly"
)

// TODO import from friendly
func writeToNewFile(fname string, data string) error {
	f, err := os.Create(fname)
	if err != nil {
		return err
	}

	w := bufio.NewWriter(f)
	_, err = w.WriteString("[]")
	if err != nil {
		return err
	}

	err = w.Flush()
	return err
}

func New(projectName string) error {
	for _, c := range []string{"/", ".", "\\"} {
		if strings.Contains(projectName, c) {
			return errors.New("invalid project name")
		}
	}

	// TODO import from friendly Exists()
	if friendly.IsDir(projectName) {
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

	err := writeToNewFile(path.Join(projectName, "res", "data", "menu.yml"), "[]")
	if err != nil {
		return err
	}

	exe, err := os.Executable()
	if err != nil {
		return err
	}

	err = writeToNewFile(path.Join(projectName, "Makefile"), fmt.Sprintf(`
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

	err = writeToNewFile(path.Join(projectName, "README.md"), "# helloworld application for the ti-nspire")
	if err != nil {
		return err
	}

	err = writeToNewFile(path.Join(projectName, "app.lua"), `
hello_world_element = Components.Base.TextField:new()
hello_world_element.Label = "Hello World"

App:AddElement(hello_world_element)
`)
	if err != nil {
		return err
	}

	return nil
}
