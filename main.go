package main

import (
	"fmt"
	"os"
)

func main() {
	switch os.Args[1] {
	case "new":
		if len(os.Args) < 3 {
			fmt.Println("You need to specify the project name")
			os.Exit(1)
		}
		err := New(os.Args[2])
		if err != nil {
			fmt.Printf("FATAL ERROR: %s", err.Error())
			os.Exit(1)
		}
	case "build":
		if len(os.Args) < 3 {
			fmt.Println("You need to specify the project folder")
			os.Exit(1)
		}
		err := Build(os.Args[2])
		if err != nil {
			fmt.Printf("FATAL ERROR: %s\n", err.Error())
			os.Exit(1)
		}
	default:
		fmt.Println("Command not known")
		os.Exit(1)
	}
}
