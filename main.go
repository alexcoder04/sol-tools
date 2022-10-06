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
		New(os.Args[2])
	default:
		fmt.Println("Command not known")
		os.Exit(1)
	}
}
