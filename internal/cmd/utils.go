package cmd

import (
	"fmt"
	"os/user"
	"time"

	"github.com/mvisonneau/go-helpers/logger"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var start time.Time

func configure(ctx *cli.Context) error {
	if !isRunningUserRoot() {
		return fmt.Errorf("You must be running as root")
	}

	start = ctx.App.Metadata["startTime"].(time.Time)

	// Configure logger
	return logger.Configure(logger.Config{
		Level:  ctx.String("log-level"),
		Format: ctx.String("log-format"),
	})
}

func isRunningUserRoot() bool {
	if user, err := user.Current(); err != nil {
		return false
	} else if user.Uid != "0" {
		return false
	}
	return true
}

func exit(exitCode int, err error) cli.ExitCoder {
	defer log.WithFields(
		log.Fields{
			"execution-time": time.Since(start),
		},
	).Debug("exited..")

	if err != nil {
		log.Error(err.Error())
	}

	return cli.NewExitError("", exitCode)
}

// ExecWrapper gracefully logs and exits our `run` functions
func ExecWrapper(f func(ctx *cli.Context) (int, error)) cli.ActionFunc {
	return func(ctx *cli.Context) error {
		return exit(f(ctx))
	}
}
