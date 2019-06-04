package cli

import (
	"time"

	"github.com/mvisonneau/automount/command"
	"github.com/urfave/cli"
)

// Init : Generates CLI configuration for the application
func Init(version *string, start time.Time) (app *cli.App) {
	app = cli.NewApp()
	app.Name = "automount"
	app.Version = *version
	app.Compiled = time.Now()
	app.Usage = "Automatically format and mount block devices"
	app.EnableBashCompletion = true

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "log-level",
			EnvVar: "AUTOMOUNT_LOG_LEVEL",
			Usage:  "log `level` (debug,info,warn,fatal,panic)",
			Value:  "info",
		},
		cli.StringFlag{
			Name:   "log-format",
			EnvVar: "AUTOMOUNT_LOG_FORMAT",
			Usage:  "log `format` (json,text)",
			Value:  "text",
		},
		cli.StringFlag{
			Name:   "device, d",
			EnvVar: "AUTOMOUNT_DEVICE",
			Usage:  "block device to mount",
			Value:  "auto",
		},
		cli.StringFlag{
			Name:   "fstype, t",
			EnvVar: "AUTOMOUNT_FSTYPE",
			Usage:  "fs type to use for the block device to mount",
			Value:  "ext4",
		},
		cli.BoolFlag{
			Name:   "use-formatted-devices",
			EnvVar: "AUTOMOUNT_USE_FORMATTED_DEVICES",
			Usage:  "use formatted but unconfigured devices (will reformat them!)",
		},
		cli.IntFlag{
			Name:   "mountpoint-mode, m",
			EnvVar: "AUTOMOUNT_MOUNTPOINT_MODE",
			Usage:  "file permissions to ensure on the mountpoint",
			Value:  0755,
		},
	}

	app.Commands = []cli.Command{
		{
			Name:      "mount",
			Usage:     "format and mount a block device somewhere",
			ArgsUsage: "<mountpoint>",
			Action:    command.Mount,
		},
		{
			Name:      "validate",
			Usage:     "check the status of dependencies",
			ArgsUsage: "",
			Action:    command.Validate,
		},
	}

	app.Metadata = map[string]interface{}{
		"startTime": start,
	}

	return
}
