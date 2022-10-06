package main

import (
	"fmt"

	"github.com/alexcoder04/friendly"
)

func main() {
	if !friendly.IsDir(".") {
		return
	}
	fmt.Println("Hello world")
}
