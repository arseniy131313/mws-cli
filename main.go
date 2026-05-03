package main

import (
	"os"

	"mws-cli/cmd/mws"
)

func main() {
	app := mws.New(os.Stdout, os.Stderr)
	os.Exit(app.Run(os.Args[1:]))
}
