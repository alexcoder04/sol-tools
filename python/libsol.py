#!/usr/bin/env python3
# Some Python functions still in use for Sol in the Go version
# TODO: only temporary solution

import sys
import yaml

def pydict_toluatable(data: dict, prefix: str) -> str:
    lines = [f"\n{prefix} = {{}}"]
    for key in data:
        res = lua_escape(data[key], prefix)
        if type(res) == list:
            for i in res:
                lines.append(i)
            continue
        lines.append(f"\n{prefix}.{key} = {res}")
    return lines

def compile_data(folder: str, filename: str) -> None:
    print(f"#compile res/data/{filename}")
    with open(f"{folder}/res/data/{filename}", "r") as fi:
        data = yaml.safe_load(fi)
        dname = filename.replace(".yml", "")
        print(f"-- BEGIN compile:{filename}")
        for line in pydict_toluatable(data, f"App.Data.Const.{dname}"):
            print(line)
        print(f"-- END compile:{filename}")

def compile_component(folder: str, filename: str) -> (str, list[str], str):
    print(f"#compile components/{filename}")
    with open(f"{folder}/components/{filename}", "r") as f:
        comp = yaml.safe_load(f)
        comp_name = file.split(".")[0]
        lua_code = [
            f"\nComponents.Custom.{comp_name} = Components.{comp['Inherit']}:new()",
            f"function Components.Custom.{comp_name}:new(o)",
            f"  o = o or Components.{comp['Inherit']}:new(o)",
            "  setmetatable(o, self)",
            "  self.__index = self"
            ]
        for key in comp:
            if key == "Inherit":
                continue
            if key in ("Update", "OnClick"):
                lua_code.append(f"  function self:{key}() {comp[key]} end")
                continue
            if key == "Color" and type(comp[key]) != list:
                lua_code.append(f"  self.{key} = {comp[key]}")
                continue
            if type(comp[key]) == str:
                value = f"\"{comp[key]}\""
            elif type(comp[key]) == bool:
                value = str(comp[key]).lower()
            elif type(comp[key]) == list:
                value = str(comp[key]).replace("[", "{").replace("]", "}")
            else:
                value = comp[key]
            lua_code.append(f"  self.{key} = {value}")
        lua_code.append("  return o")
        lua_code.append("end")
    return comp_name, lua_code, comp["Inherit"]

def compile_components(folder: str) -> None:
        components = []
        for file in os.listdir(f"{folder}/components"):
            components.append(compile_component(file, fo))
        components_sorted = []
        loop_detector = 0
        while len(components) > 0:
            delete = []
            for c in components:
                if c[2].startswith("Base."):
                    components_sorted.append(c)
                    delete.append(c)
                    continue
                if c[2] in [f"Custom.{j[0]}" for j in components_sorted]:
                    components_sorted.append(c)
                    delete.append(c)
                    continue
            for c in delete:
                components.remove(c)
            loop_detector += 1
            if loop_detector > 99:
                die("Component inheritance loop deteted")
        for i in components_sorted:
            print(f"\n-- BEGIN compile:components/{i[0]}")
            for line in i[1]:
                print(line + "\n")
            print(f"-- END compile:components/{i[0]}")

def compile_menu(folder: str) -> None:
    with open(f"{folder}/res/data/menu.yml", "r") as fi:
        menu = yaml.safe_load(fi)
        menu.append({
            "Id": "help",
            "Name": "Help",
            "Submenues": [
                {
                    "Id": "about",
                    "Name": "About",
                    "Function": "Library.Internal:ShowAboutDialog()"
                }
            ]
        })
        categories = []
        functions = []
        for cat in menu:
            submenues = []
            for sm in cat["Submenues"]:
                submenues.append(f"{{\"{sm['Name']}\", _menu_{cat['Id']}_{sm['Id']}}}")
            functions.append(f"function _menu_{cat['Id']}_{sm['Id']}() {sm['Function']} end")
            categories.append(f"{{\"{cat['Name']}\", {', '.join(submenues)}}}")
    print(f"-- BEGIN compile:res/data/menu.yml")
    for i in functions:
        print(i)
    print(f"toolpalette.register({{{', '.join(categories)}}})")
    print(f"-- END compile:res/data/menu.yml")

if __name__ == "__main__":
    if len(sys.argv) < 4:
        print("too few arguments")
        sys.exit(1)
    
    if sys.argv[1] == "compile_data":
        compile_data(sys.argv[2], sys.argv[3])

    if sys.argv[1] == "compile_components":
        compile_components(sys.argv[2])

    if sys.argv[1] == "compile_menu":
        compile_menu(sys.argv[2])
