package utils

import (
	"bufio"
	"fmt"
	"os"
)

func Exit(code int, message string) {
	if message != "" {
		fmt.Println(message)
	}
	fmt.Print("Sol is done. Press enter to exit")
	r := bufio.NewReader(os.Stdin)
	_, _ = r.ReadString('\n')
	os.Exit(code)
}
