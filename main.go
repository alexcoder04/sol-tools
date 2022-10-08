package main

import (
	"os"

	"github.com/alexcoder04/arrowprint"
)

func main() {
	switch os.Args[1] {
	case "new":
		if len(os.Args) < 3 {
			arrowprint.Err0("You need to specify the project name")
			os.Exit(1)
		}
		err := New(os.Args[2])
		if err != nil {
			arrowprint.Err0("Something went wrong when creating your project: %s", err.Error())
			os.Exit(1)
		}
	case "build":
		if len(os.Args) < 3 {
			arrowprint.Err0("You need to specify the project folder")
			os.Exit(1)
		}
		err := Build(os.Args[2])
		if err != nil {
			arrowprint.Err0("Something went wrong when building your project: %s\n", err.Error())
			os.Exit(1)
		}
	default:
		arrowprint.Err0("Command not known")
		os.Exit(1)
	}
}
