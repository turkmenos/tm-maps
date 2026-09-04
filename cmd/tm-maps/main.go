package main

import (
	"os"

	"github.com/turkmenos/tm-maps/internal/cli"
)

func main() {
	os.Exit(cli.Run(os.Stdout, os.Stderr, os.Args[1:]))
}
