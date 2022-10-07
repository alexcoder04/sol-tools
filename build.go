package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"

	"github.com/alexcoder04/friendly"
)

var OUT_LUA = path.Join(os.TempDir(), "out.lua")

// TODO find actual library path
const LUALIB = "/home/alex/Repos/sol-lib"

// TODO replace with go code
const PYTHONLIB = "/home/alex/Repos/sol-tools/python/libsol.py"

func appendFile(projectFolder string, filename string, type_ string, w *bufio.Writer) error {
	var prefix string = "."
	if type_ == "include" {
		prefix = LUALIB
	}
	if type_ == "process" {
		prefix = projectFolder
	}

	fmt.Printf("#%s %s\n", type_, filename)

	content, err := ioutil.ReadFile(path.Join(prefix, filename))
	if err != nil {
		return err
	}

	_, err = w.WriteString(fmt.Sprintf("\n-- BEGIN %s:%s", type_, filename))
	if err != nil {
		return err
	}
	if filename == "init.lua" && type_ == "process" {
		_, err := w.WriteString("\nfunction init()")
		if err != nil {
			return err
		}
	}

	_, err = w.Write(content)
	if err != nil {
		return err
	}

	if filename == "init.lua" && type_ == "process" {
		_, err := w.WriteString("\nend")
		if err != nil {
			return err
		}
	}
	_, err = w.WriteString(fmt.Sprintf("\n-- END %s:%s", type_, filename))
	return err
}

func appendFolder(projectFolder string, subfolder string, type_ string, w *bufio.Writer) error {
	err := appendFile(projectFolder, path.Join(subfolder, "_init.lua"), type_, w)
	if err != nil {
		return err
	}

	var folderAbs string = "."
	if type_ == "include" {
		folderAbs = path.Join(LUALIB, subfolder)
	}
	if type_ == "process" {
		folderAbs = path.Join(projectFolder, subfolder)
	}

	files, err := ioutil.ReadDir(folderAbs)
	if err != nil {
		return err
	}

	for _, file := range files {
		if file.Name() == "_init.lua" {
			continue
		}
		err := appendFile(projectFolder, path.Join(subfolder, file.Name()), type_, w)
		if err != nil {
			return err
		}
	}

	return nil
}

// TODO import from friendly
func prepareCommand(command string, arguments []string, workingDir string) *exec.Cmd {
	if workingDir == "" {
		workingDir = friendly.Getpwd()
	}

	cmd := exec.Command(command, arguments...)
	cmd.Dir = workingDir

	return cmd
}

// TODO import from friendly
func GetOutput(command string, arguments []string, workingDir string) (string, error) {
	cmd := prepareCommand(command, arguments, workingDir)
	out, err := cmd.Output()
	return string(out), err
}

func appendData(projectFolder string, w *bufio.Writer) error {
	files, err := ioutil.ReadDir(path.Join(projectFolder, "res", "data"))
	if err != nil {
		return err
	}
	for _, file := range files {
		fmt.Printf("#compile res/data/%s -> skipping!\n", file.Name())
		// TODO
	}
	return nil
}

func appendComponents(projectFolder string, w *bufio.Writer) error {
	fmt.Println("#compile components/")
	data, err := GetOutput(PYTHONLIB, []string{
		"compile_components",
		projectFolder,
		""}, "")
	if err != nil {
		return err
	}
	_, err = w.WriteString(data)
	return err
}

func appendMenu(projectFolder string, w *bufio.Writer) error {
	fmt.Println("#compile res/data/menu.yml")
	data, err := GetOutput(PYTHONLIB, []string{
		"compile_menu",
		projectFolder,
		""}, "")
	if err != nil {
		return err
	}
	_, err = w.WriteString(data)
	return err
}

func resolvePath(folder string) (string, error) {
	if path.IsAbs(folder) {
		return folder, nil
	}
	return filepath.Abs(folder)
}

func Build(pfolder string) error {
	// TODO use Exists() from friendly
	if friendly.IsFile(OUT_LUA) {
		err := os.RemoveAll(OUT_LUA)
		if err != nil {
			return err
		}
	}

	folder, err := resolvePath(pfolder)
	if err != nil {
		return err
	}

	if !friendly.IsDir(folder) {
		return os.ErrNotExist
	}

	f, err := os.Create(OUT_LUA)
	if err != nil {
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)

	// framework: main file
	err = appendFile(folder, "app.lua", "include", w)
	if err != nil {
		return err
	}
	// framework: library
	err = appendFolder(folder, "library", "include", w)
	if err != nil {
		return err
	}
	// framework: components
	err = appendFolder(folder, "components", "include", w)
	if err != nil {
		return err
	}

	// project: data
	err = appendData(folder, w)
	if err != nil {
		return err
	}
	// project: components
	err = appendComponents(folder, w)
	if err != nil {
		return err
	}
	// project: lua code
	for _, file := range []string{"app.lua", "init.lua", "hooks.lua"} {
		err := appendFile(folder, file, "process", w)
		if err != nil {
			return err
		}
	}
	// project: menu
	err = appendMenu(folder, w)
	if err != nil {
		return err
	}

	// library: events
	err = appendFile(folder, "events.lua", "include", w)
	if err != nil {
		return err
	}

	err = w.Flush()
	if err != nil {
		return err
	}

	fmt.Println("All files merged succesfully")
	return nil
}
