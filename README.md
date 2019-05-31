# automount

[![GoDoc](https://godoc.org/github.com/mvisonneau/automount?status.svg)](https://godoc.org/github.com/mvisonneau/automount)
[![Go Report Card](https://goreportcard.com/badge/github.com/mvisonneau/automount)](https://goreportcard.com/report/github.com/mvisonneau/automount)
[![Docker Pulls](https://img.shields.io/docker/pulls/mvisonneau/automount.svg)](https://hub.docker.com/r/mvisonneau/automount/)
[![Build Status](https://cloud.drone.io/api/badges/mvisonneau/automount/status.svg)](https://cloud.drone.io/mvisonneau/automount)
[![Coverage Status](https://coveralls.io/repos/github/mvisonneau/automount/badge.svg?branch=master)](https://coveralls.io/github/mvisonneau/automount?branch=master)

Automatically mount block devices

## TL;DR

Assuming you have a drive available (allocatable & not formatted, nor mounted) on your machine :

```bash
~$ lsblk
NAME        MAJ:MIN RM   SIZE RO TYPE MOUNTPOINT
nvme0n1     259:0    0 372.5G  0 disk
nvme1n1     259:1    0    10G  0 disk
└─nvme1n1p1 259:2    0    10G  0 part /

# Automatically format and mount it where you need it.
~$ automount mount /mnt/foo
INFO[2019-05-23T08:21:14Z] No device specified, trying to find one automatically
INFO[2019-05-23T08:21:14Z] Found 2 disk(s), total size of 382 GB
INFO[2019-05-23T08:21:14Z] /dev/nvme0n1 is available, picking this one!
INFO[2019-05-23T08:21:14Z] Formatting device /dev/nvme0n1 to ext4
INFO[2019-05-23T08:21:17Z] Parsing /etc/fstab
INFO[2019-05-23T08:21:17Z] Found 2 entries in /etc/fstab
INFO[2019-05-23T08:21:17Z] Ensuring that mountpoint /mnt/foo exists with correct permissions (493)
INFO[2019-05-23T08:21:17Z] /dev/nvme0n1 is not configured within fstab, appending configuration
INFO[2019-05-23T08:21:17Z] Writing configuration to /etc/fstab
INFO[2019-05-23T08:21:17Z] Attempting to mount /dev/nvme0n1 to /mnt/foo
INFO[2019-05-23T08:21:17Z] Mounted!

~$ df -h /mnt/foo
Filesystem      Size  Used Avail Use% Mounted on
/dev/nvme0n1    366G   69M  347G   1% /mnt/foo
```

## Usage

```
~$ automount
NAME:
   automount - Automatically format and mount block devices

USAGE:
   automount [global options] command [command options] [arguments...]

COMMANDS:
     mount     format and mount a block device somewhere
     validate  check the status of dependencies
     help, h   Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --log-level level                  log level (debug,info,warn,fatal,panic) (default: "info") [$AUTOMOUNT_LOG_LEVEL]
   --log-format format                log format (json,text) (default: "text") [$AUTOMOUNT_LOG_FORMAT]
   --device value, -d value           block device to mount (default: "auto") [$AUTOMOUNT_DEVICE]
   --fstype value, -t value           fs type to use for the block device to mount (default: "ext4") [$AUTOMOUNT_FSTYPE]
   --mountpoint-mode value, -m value  file permissions to ensure on the mountpoint (default: 493) [$AUTOMOUNT_MOUNTPOINT_MODE]
   --help, -h                         show help
   --version, -v                      print the version
```

## Prerequisites

- Linux OS
- `blkid`, `lsblk`, `mkfs.*` commands available

## Contribute

Contributions are more than welcome! Feel free to submit a [PR](https://github.com/mvisonneau/automount/pulls).
