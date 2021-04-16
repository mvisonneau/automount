package cli

import (
	"fmt"
	"os"
	"time"

	"github.com/mvisonneau/automount/internal/cmd"
	"github.com/urfave/cli/v2"
)

// Run handles the instanciation of the CLI application
func Run(version string, args []string) {
	err := NewApp(version, time.Now()).Run(args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// NewApp configures the CLI application
func NewApp(version string, start time.Time) (app *cli.App) {
	app = cli.NewApp()
	app.Name = "automount"
	app.Version = version
	app.Usage = "Automatically format and mount block devices"
	app.EnableBashCompletion = true

	app.Flags = cli.FlagsByName{
		&cli.StringFlag{
			Name:    "log-level",
			EnvVars: []string{"AUTOMOUNT_LOG_LEVEL"},
			Usage:   "log `level` (trace,debug,info,warn,fatal,panic)",
			Value:   "info",
		},
		&cli.StringFlag{
			Name:    "log-format",
			EnvVars: []string{"AUTOMOUNT_LOG_FORMAT"},
			Usage:   "log `format` (json,text)",
			Value:   "text",
		},
		&cli.StringSliceFlag{
			Name:    "device",
			Aliases: []string{"d"},
			EnvVars: []string{"AUTOMOUNT_DEVICES"},
			Usage:   "block device(s) to mount",
		},
		&cli.StringFlag{
			Name:    "fstype",
			Aliases: []string{"t"},
			EnvVars: []string{"AUTOMOUNT_FSTYPE"},
			Usage:   "fs type to use for the block device to mount",
			Value:   "ext4",
		},
		&cli.BoolFlag{
			Name:    "use-formatted-devices",
			EnvVars: []string{"AUTOMOUNT_USE_FORMATTED_DEVICES"},
			Usage:   "use formatted but unconfigured devices (will reformat them!)",
		},
		&cli.BoolFlag{
			Name:    "use-lvm",
			EnvVars: []string{"AUTOMOUNT_USE_LVM"},
			Usage:   "use LVM for the partitioning of the block devices",
		},
		&cli.BoolFlag{
			Name:    "use-all-devices",
			EnvVars: []string{"AUTOMOUNT_USE_ALL_DEVICES"},
			Usage:   "use all available devices in a soft-raid fashion (requires --use-lvm as well)",
		},
		&cli.IntFlag{
			Name:    "mountpoint-mode",
			Aliases: []string{"m"},
			EnvVars: []string{"AUTOMOUNT_MOUNTPOINT_MODE"},
			Usage:   "file permissions to ensure on the mountpoint",
			Value:   0o755,
		},
	}

	app.Commands = cli.CommandsByName{
		&cli.Command{
			Name:      "mount",
			Usage:     "format and mount a block device somewhere",
			ArgsUsage: "<mountpoint>",
			Action:    cmd.ExecWrapper(cmd.Mount),
		},
		&cli.Command{
			Name:      "validate",
			Usage:     "check the status of dependencies",
			ArgsUsage: "",
			Action:    cmd.ExecWrapper(cmd.Validate),
		},
	}

	app.Metadata = map[string]interface{}{
		"startTime": start,
	}

	return
}
