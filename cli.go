package main

import (
	"github.com/urfave/cli"
)

var version = "<devel>"

// runCli : Generates cli configuration for the application
func runCli() (c *cli.App) {
	c = cli.NewApp()
	c.Name = "automount"
	c.Version = version
	c.Usage = "Automatically format and mount block devices"
	c.EnableBashCompletion = true

	c.Flags = []cli.Flag{
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
		cli.IntFlag{
			Name:   "mountpoint-mode, m",
			EnvVar: "AUTOMOUNT_MOUNTPOINT_MODE",
			Usage:  "file permissions to ensure on the mountpoint",
			Value:  0755,
		},
	}

	c.Commands = []cli.Command{
		{
			Name:      "mount",
			Usage:     "format and mount a block device somewhere",
			ArgsUsage: "<mountpoint>",
			Action:    executeMount,
		},
	}

	return
}
