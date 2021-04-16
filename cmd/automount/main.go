package main

import (
	"os"

	"github.com/mvisonneau/automount/internal/cli"
)

var version = ""

func main() {
	cli.Run(version, os.Args)
}
