package main

import (
	"log"
	"os"
	"time"

	"github.com/mvisonneau/automount/cli"
)

var version = ""

func main() {
	if err := cli.Init(&version, time.Now()).Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
