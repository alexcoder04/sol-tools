package main

import (
	"flag"
	"strings"

	"github.com/alexcoder04/arrowprint"
	"github.com/alexcoder04/friendly/v2"
	"github.com/alexcoder04/sol-tools/commands"
	"github.com/alexcoder04/sol-tools/utils"
)

var AdvancedMode = flag.Bool("a", false, "use advanced/automated mode")

func Run(args []string) {
	switch args[0] {
	case "new":
		if len(args) < 3 {
			arrowprint.Err0("You need to specify the project name and parent folder")
			utils.Exit(1, "", *AdvancedMode)
		}
		err := commands.New(args[1], args[2])
		if err != nil {
			arrowprint.Err0("Something went wrong when creating your project: %s", err.Error())
			utils.Exit(1, "", *AdvancedMode)
		}
	case "build":
		if len(args) < 2 {
			arrowprint.Err0("You need to specify the project folder")
			utils.Exit(1, "", *AdvancedMode)
		}
		err := commands.Build(args[1])
		if err != nil {
			arrowprint.Err0("Something went wrong when building your project: %s\n", err.Error())
			utils.Exit(1, "", *AdvancedMode)
		}
	case "help":
		commands.Help()
		utils.Exit(0, "", *AdvancedMode)
	default:
		arrowprint.Err0("Unknown command: '%s'", args[0])
		utils.Exit(1, "", *AdvancedMode)
	}
	utils.Exit(0, "", *AdvancedMode)
}

func main() {
	flag.Parse()

	if *AdvancedMode {
		Run(flag.Args())
	}

	cmd, err := friendly.Input("Type a Command: ")
	if err != nil {
		utils.Exit(1, "Error reading command", false)
	}
	cmd = strings.TrimSpace(cmd)
	switch cmd {
	case "new":
		folder, err := friendly.Input("Type or drag the folder to create the project in: ")
		if err != nil {
			utils.Exit(1, "Error reading folder", false)
		}
		name, err := friendly.Input("Type the name of your project: ")
		if err != nil {
			utils.Exit(1, "Error reading name", false)
		}
		Run([]string{cmd, strings.TrimSpace(folder), strings.TrimSpace(name)})
	case "build":
		folder, err := friendly.Input("Type or drag the folder to build: ")
		if err != nil {
			utils.Exit(1, "Error reading folder", false)
		}
		Run([]string{cmd, strings.TrimSpace(folder)})
	case "help":
		Run([]string{cmd})
	default:
		arrowprint.Err0("Unknown command: '%s'", cmd)
		arrowprint.Info1("Type 'help' for a list of commands")
		utils.Exit(1, "", false)
	}
}
