#!/bin/sh

cd "$HOME/Repos/sol-tools"
go build -o sol .
cd "$OLDPWD"
"$HOME/Repos/sol-tools/sol" "$@"
