package main

import (
	"os"

	"lotka-volt/internal/cli"
)

func main() {
	os.Exit(cli.Run(os.Args[1:]))
}
