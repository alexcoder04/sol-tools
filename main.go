package main

import (
	"flag"

	"github.com/alexcoder04/arrowprint"
)

var AdvancedMode = flag.Bool("a", false, "use advanced/automated mode")

func Run(args []string) {
	switch args[0] {
	case "new":
		if len(args) < 2 {
			arrowprint.Err0("You need to specify the project name")
			Exit(1, "")
		}
		err := New(args[1])
		if err != nil {
			arrowprint.Err0("Something went wrong when creating your project: %s", err.Error())
			Exit(1, "")
		}
	case "build":
		if len(args) < 2 {
			arrowprint.Err0("You need to specify the project folder")
			Exit(1, "")
		}
		err := Build(args[1])
		if err != nil {
			arrowprint.Err0("Something went wrong when building your project: %s\n", err.Error())
			Exit(1, "")
		}
	default:
		arrowprint.Err0("Command '%s' not known", args[0])
		Exit(1, "")
	}
	Exit(0, "")
}

func main() {
	flag.Parse()

	if *AdvancedMode {
		Run(flag.Args())
	}

	cmd, err := Input("Type a Command: ")
	if err != nil {
		Exit(1, "Error reading command")
	}
	switch cmd {
	case "new":
		folder, err := Input("Type or drag the folder to create the project in")
		if err != nil {
			Exit(1, "Error reading folder")
		}
		Run([]string{cmd, folder})
	case "build":
		folder, err := Input("Type or drag the folder to build")
		if err != nil {
			Exit(1, "Error reading folder")
		}
		Run([]string{cmd, folder})
	}
}
