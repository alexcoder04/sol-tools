#!/usr/bin/env python3
# Some Python functions still in use for Sol in the Go version
# TODO: only temporary solution

import os
import sys
import yaml

def compile_component(folder: str, filename: str) -> (str, list[str], str):
    with open(f"{folder}/components/{filename}", "r") as f:
        comp = yaml.safe_load(f)
        comp_name = filename.split(".")[0]
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
            components.append(compile_component(folder, file))
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
            print(f"-- BEGIN compile:components/{i[0]}")
            for line in i[1]:
                print(line + "\n")
            print(f"-- END compile:components/{i[0]}")

if __name__ == "__main__":
    if len(sys.argv) < 4:
        print("too few arguments")
        sys.exit(1)
    
    if sys.argv[1] == "compile_components":
        compile_components(sys.argv[2])
