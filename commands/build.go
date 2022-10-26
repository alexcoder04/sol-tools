package commands

import (
	"bufio"
	"os"
	"path"
	"path/filepath"

	"github.com/alexcoder04/arrowprint"
	"github.com/alexcoder04/friendly"
	"github.com/alexcoder04/sol-tools/utils"
)

func Build(pfolder string) error {
	arrowprint.InfoC("Building your project at %s", pfolder)
	arrowprint.Info1("Doing initial checks")

	folder, err := filepath.Abs(pfolder)
	if err != nil {
		return err
	}

	if !friendly.IsDir(folder) {
		return os.ErrNotExist
	}

	OUT_LUA := path.Join(folder, "out.lua")

	if friendly.Exists(OUT_LUA) {
		err := os.RemoveAll(OUT_LUA)
		if err != nil {
			return err
		}
	}

	arrowprint.Info1("Reading project metadata")
	metadata, err := getMetadata(folder)
	if err != nil {
		return err
	}

	arrowprint.Info1("Getting sol library")
	libPath, version, err := utils.GetLibrary(metadata["SolVersion"])
	if err != nil {
		return err
	}
	metadata["SolVersion"] = version

	f, err := os.Create(OUT_LUA)
	if err != nil {
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)

	// framework: main file
	err = appendFile(folder, libPath, "app.lua", "include", w)
	if err != nil {
		return err
	}
	// framework: library
	err = appendFolder(folder, libPath, "library", "include", w)
	if err != nil {
		return err
	}
	// framework: components
	err = appendFolder(folder, libPath, "components", "include", w)
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
	// project: metadata
	err = appendMetadata(metadata, w)
	if err != nil {
		return err
	}
	// project: lua code
	for _, file := range []string{"app.lua", "init.lua", "hooks.lua"} {
		if !friendly.IsFile(path.Join(folder, file)) && file != "app.lua" {
			continue
		}
		err := appendFile(folder, libPath, file, "process", w)
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
	err = appendFile(folder, libPath, "events.lua", "include", w)
	if err != nil {
		return err
	}

	err = w.Flush()
	if err != nil {
		return err
	}

	arrowprint.Suc0("All files merged succesfully")
	arrowprint.Suc1("Result written to %s", OUT_LUA)

	arrowprint.Info0("You can now either copy the result into the TI student software or compile it into .tns with Luna")
	return nil
}
