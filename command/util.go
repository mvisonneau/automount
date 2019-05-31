package command

import (
	"fmt"
	"os/user"
	"time"

	"github.com/mvisonneau/automount/logger"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var start time.Time

func configure(ctx *cli.Context) error {
	if !isRunningUserRoot() {
		return fmt.Errorf("You must be running as root")
	}

	start = ctx.App.Metadata["startTime"].(time.Time)

	l := &logger.Logger{
		Level:  ctx.GlobalString("log-level"),
		Format: ctx.GlobalString("log-format"),
	}

	return l.Configure()
}

func isRunningUserRoot() bool {
	if user, err := user.Current(); err != nil {
		return false
	} else if user.Uid != "0" {
		return false
	}
	return true
}

func exit(err error, exitCode int) *cli.ExitError {
	defer log.Debugf("Executed in %s, exiting..", time.Since(start))
	if err != nil {
		log.Error(err.Error())
		return cli.NewExitError("", exitCode)
	}

	return nil
}
