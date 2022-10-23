package commands

import "fmt"

func Help() {
	fmt.Println(`
Commands:
	build		compile a sol project for uploading to calculator
	help		show this help
	new			create a new sol project
	`)
}
