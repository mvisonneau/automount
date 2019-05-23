package main

import (
	"os"
	"time"
)

var start time.Time

func main() {
	start = time.Now()
	runCli().Run(os.Args)
}
