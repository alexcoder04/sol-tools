package main

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"path"
	"strings"

	"gopkg.in/yaml.v3"
)

func appendFile(projectFolder string, libFolder string, filename string, type_ string, w *bufio.Writer) error {
	var prefix string = "."
	if type_ == "include" {
		prefix = libFolder
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
	_, err = w.WriteString(fmt.Sprintf("\n-- END %s:%s\n", type_, filename))
	return err
}

func appendFolder(projectFolder string, libFolder string, subfolder string, type_ string, w *bufio.Writer) error {
	err := appendFile(projectFolder, libFolder, path.Join(subfolder, "_init.lua"), type_, w)
	if err != nil {
		return err
	}

	var folderAbs string = "."
	if type_ == "include" {
		folderAbs = path.Join(libFolder, subfolder)
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
		err := appendFile(projectFolder, libFolder, path.Join(subfolder, file.Name()), type_, w)
		if err != nil {
			return err
		}
	}

	return nil
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

func compileComponent(projectFolder string, name string, w *bufio.Writer) (string, string, []string, error) {
	data, err := ioutil.ReadFile(path.Join(projectFolder, "components", name))
	if err != nil {
		return "", "", []string{}, err
	}

	comp := make(map[string]string)
	err = yaml.Unmarshal(data, &comp)
	if err != nil {
		return "", "", []string{}, err
	}
	compName := strings.TrimRight(name, ".yml")
	luaCode := []string{
		fmt.Sprintf("\nComponents.Custom.%s = Components.%s:new()", compName, comp["Inherit"]),
		fmt.Sprintf("function Components.Custom.%s:new(o)", compName),
		fmt.Sprintf("  o = o or Components.%s:new(o)", comp["Inherit"]),
		"  setmetatable(o, self)",
		"  self.__index = self"}

	for key, val := range comp {
		if key == "Inherit" {
			continue
		}
		if key == "OnClick" {
			luaCode = append(luaCode, fmt.Sprintf("  function self:%s() %s end", key, val))
			continue
		}
		if key == "Update" {
			luaCode = append(luaCode, fmt.Sprintf("  function self:%s(focused) %s end", key, val))
			continue
		}
		if key == "Color" && strings.HasPrefix(val, "Lib.Colors.") {
			luaCode = append(luaCode, fmt.Sprintf("  self.%s = %s", key, val))
			continue
		}
		if val == "true" || val == "false" || isNumber(val) {
			luaCode = append(luaCode, fmt.Sprintf("  self.%s = %s", key, val))
			continue
		}
		if strings.HasPrefix(val, "[") && strings.HasSuffix(val, "]") {
			luaCode = append(luaCode, fmt.Sprintf("  self.%s = {%s}", key, val[1:len(val)-2]))
			continue
		}
		luaCode = append(luaCode, fmt.Sprintf("  self.%s = \"%s\"", key, val))
	}

	luaCode = append(luaCode, "  return o")
	luaCode = append(luaCode, "end")
	return compName, comp["Inherit"], luaCode, nil
}

func appendComponents(projectFolder string, w *bufio.Writer) error {
	fmt.Println("#compile components/")
	components := make(map[int]Component)
	files, err := ioutil.ReadDir(path.Join(projectFolder, "components"))
	if err != nil {
		return err
	}

	for i, f := range files {
		comp := Component{}
		n, p, c, err := compileComponent(projectFolder, f.Name(), w)
		if err != nil {
			return err
		}
		comp.Name = n
		comp.Parent = p
		comp.Code = c
		components[i] = comp
	}

	var componentsSorted []Component
	loopDetector := 0
	for len(components) > 0 {
		var forDelete []int
		for i, c := range components {
			if strings.HasPrefix(c.Parent, "Base.") {
				componentsSorted = append(componentsSorted, c)
				forDelete = append(forDelete, i)
				continue
			}
			availableComps := []string{}
			for _, ac := range componentsSorted {
				availableComps = append(availableComps, "Custom."+ac.Name)
			}
			if stringArrayContains(availableComps, c.Parent) {
				componentsSorted = append(componentsSorted, c)
				forDelete = append(forDelete, i)
				continue
			}
		}
		for _, i := range forDelete {
			delete(components, i)
		}
		loopDetector += 1
		if loopDetector > 99 {
			return errors.New("Component inheritance loop detected")
		}
	}

	for _, c := range componentsSorted {
		_, err := w.WriteString(fmt.Sprintf("\n-- BEGIN compile:components/%s\n", c.Name))
		if err != nil {
			return err
		}
		for _, line := range c.Code {
			_, err := w.WriteString(line + "\n")
			if err != nil {
				return err
			}
		}
		_, err = w.WriteString(fmt.Sprintf("\n-- END compile:components/%s\n", c.Name))
		if err != nil {
			return err
		}
	}
	return nil
}

func appendMenu(projectFolder string, w *bufio.Writer) error {
	fmt.Println("#compile res/data/menu.yml")
	var menu []MenuEntry
	data, err := ioutil.ReadFile(path.Join(projectFolder, "res", "data", "menu.yml"))
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(data, &menu)
	if err != nil {
		return err
	}
	menu = append(menu, MenuEntry{
		"help",
		"Help",
		[]SubmenuEntry{
			{
				"about",
				"About",
				"Lib.Internal.ShowAboutDialog()"}}})

	var categories []string
	var functions []string
	for _, c := range menu {
		var submenues []string
		for _, s := range c.Submenues {
			submenues = append(submenues, fmt.Sprintf(
				"{\"%s\", _menu_%s_%s}",
				s.Name, c.Id, s.Id))
			functions = append(functions, fmt.Sprintf(
				"function _menu_%s_%s() %s end\n",
				c.Id, s.Id, s.Function))
		}
		categories = append(categories, fmt.Sprintf(
			"{\"%s\", %s}",
			c.Name, strings.Join(submenues, ", ")))
	}

	_, err = w.WriteString("\n-- BEGIN compile:res/data/menu.yml\n")
	if err != nil {
		return err
	}
	for _, f := range functions {
		_, err := w.WriteString(f)
		if err != nil {
			return err
		}
	}
	_, err = w.WriteString(fmt.Sprintf(
		"toolpalette.register({%s})",
		strings.Join(categories, ", ")))
	if err != nil {
		return err
	}
	_, err = w.WriteString("\n-- END compile:res/data/menu.yml\n")
	return err
}