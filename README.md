# 🗻 automount

[![PkgGoDev](https://pkg.go.dev/badge/github.com/mvisonneau/automount)](https://pkg.go.dev/mod/github.com/mvisonneau/automount)
[![Go Report Card](https://goreportcard.com/badge/github.com/mvisonneau/automount)](https://goreportcard.com/report/github.com/mvisonneau/automount)
[![Docker Pulls](https://img.shields.io/docker/pulls/mvisonneau/automount.svg)](https://hub.docker.com/r/mvisonneau/automount/)
[![test](https://github.com/mvisonneau/automount/actions/workflows/test.yml/badge.svg)](https://github.com/mvisonneau/automount/actions/workflows/test.yml)
[![Coverage Status](https://coveralls.io/repos/github/mvisonneau/automount/badge.svg?branch=main)](https://coveralls.io/github/mvisonneau/automount?branch=main)
[![release](https://github.com/mvisonneau/automount/actions/workflows/release.yml/badge.svg)](https://github.com/mvisonneau/automount/actions/workflows/release.yml)

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
~$ sudo automount mount /mnt/foo
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

## It also supports LVM and has soft-raid 0 capabilities

```bash
~$ sudo automount --log-level debug--use-lvm --use-all-devices mount /mnt/foo
INFO[2019-06-06T17:28:21Z] Parsing /etc/fstab
INFO[2019-06-06T17:28:21Z] Found 3 entries in /etc/fstab
INFO[2019-06-06T17:28:21Z] No device specified, looking up available ones
INFO[2019-06-06T17:28:21Z] Found 3 disk(s), total size of 315 GB
INFO[2019-06-06T17:28:21Z] /dev/xvda has partitions, skipping
INFO[2019-06-06T17:28:21Z] /dev/xvdb is available
INFO[2019-06-06T17:28:21Z] /dev/xvdc is available
INFO[2019-06-06T17:28:21Z] Using LVM for managing the partitions
DEBU[2019-06-06T17:28:21Z] command lvm available at /sbin/lvm
DEBU[2019-06-06T17:28:21Z] LVM: getting current state
DEBU[2019-06-06T17:28:21Z] LVM: creating physical volume on /dev/xvdb
DEBU[2019-06-06T17:28:21Z] LVM: creating physical volume on /dev/xvdc
DEBU[2019-06-06T17:28:22Z] LVM: creating volume group
DEBU[2019-06-06T17:28:23Z] LVM: creating logical volume
INFO[2019-06-06T17:28:23Z] physical volume, volume group and logical volume created, using this as a device
INFO[2019-06-06T17:28:23Z] /dev/automount/automount is not formatted, will format it.
INFO[2019-06-06T17:28:23Z] Formatting device /dev/automount/automount to ext4
INFO[2019-06-06T17:28:30Z] Ensuring that mountpoint /mnt/foo exists with correct permissions (493)
INFO[2019-06-06T17:28:30Z] /dev/automount/automount is not configured within fstab, appending configuration
INFO[2019-06-06T17:28:30Z] Writing configuration to /etc/fstab
INFO[2019-06-06T17:28:30Z] Attempting to mount /dev/automount/automount to /mnt/foo
INFO[2019-06-06T17:28:30Z] Mounted!
DEBU[2019-06-06T17:28:30Z] Executed in 9.322286497s, exiting..

~$ pvs --noheadings
  /dev/xvdb  automount lvm2 a--   152.57g    0
  /dev/xvdc  automount lvm2 a--   152.57g    0
~$ lvs --noheadings
  automount automount -wi-ao---- 305.13g
```

## Install

### Go

```bash
# >= 1.16
~$ go install github.com/mvisonneau/automount/cmd/automount@latest
```

### Docker

```bash
~$ docker run -it --rm docker.io/mvisonneau/automount
or
~$ docker run -it --rm ghcr.io/mvisonneau/automount
```

### Scoop

```bash
~$ scoop bucket add https://github.com/mvisonneau/scoops
~$ scoop install automount
```

### Binaries, DEB and RPM packages

Have a look onto the [latest release page](https://github.com/mvisonneau/automount/releases/latest) to pick your flavor and version. Here is an helper to fetch the most recent one:

```bash
~$ export AUTOMOUNT_VERSION=$(curl -s "https://api.github.com/repos/mvisonneau/automount/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
```

```bash
# Binary (eg: linux/amd64)
~$ wget https://github.com/mvisonneau/automount/releases/download/${AUTOMOUNT_VERSION}/automount_${AUTOMOUNT_VERSION}_linux_amd64.tar.gz
~$ tar zxvf automount_${AUTOMOUNT_VERSION}_linux_amd64.tar.gz -C /usr/local/bin

# DEB package (eg: linux/386)
~$ wget https://github.com/mvisonneau/automount/releases/download/${AUTOMOUNT_VERSION}/automount_${AUTOMOUNT_VERSION}_linux_386.deb
~$ dpkg -i automount_${AUTOMOUNT_VERSION}_linux_386.deb

# RPM package (eg: linux/arm64)
~$ wget https://github.com/mvisonneau/automount/releases/download/${AUTOMOUNT_VERSION}/automount_${AUTOMOUNT_VERSION}_linux_arm64.rpm
~$ rpm -ivh automount_${AUTOMOUNT_VERSION}_linux_arm64.rpm
```

## Usage

```bash
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
   --log-level level                  log level (trace,debug,info,warn,fatal,panic) (default: "info") [$AUTOMOUNT_LOG_LEVEL]
   --log-format format                log format (json,text) (default: "text") [$AUTOMOUNT_LOG_FORMAT]
   --device value, -d value           block device(s) to mount [$AUTOMOUNT_DEVICES]
   --fstype value, -t value           fs type to use for the block device to mount (default: "ext4") [$AUTOMOUNT_FSTYPE]
   --use-formatted-devices            use formatted but unconfigured devices (will reformat them!) (default: false) [$AUTOMOUNT_USE_FORMATTED_DEVICES]
   --use-lvm                          use LVM for the partitioning of the block devices (default: false) [$AUTOMOUNT_USE_LVM]
   --use-all-devices                  use all available devices in a soft-raid fashion (requires --use-lvm as well) (default: false) [$AUTOMOUNT_USE_ALL_DEVICES]
   --mountpoint-mode value, -m value  file permissions to ensure on the mountpoint (default: 493) [$AUTOMOUNT_MOUNTPOINT_MODE]
   --help, -h                         show help (default: false)
```

## Prerequisites

At the moment `automount` only works as `root` on linux based operating systems.

In order to list the dependencies, you can run the following command

```bash
~$ sudo automount validate
  COMMAND | MANDATORY | AVAILABLE
+---------+-----------+-----------+
  blkid   | YES       | YES
  lsblk   | YES       | YES
  lvm     | NO        | YES
  mdadm   | NO        | YES
```

## Contribute

Contributions are more than welcome! Feel free to submit a [PR](https://github.com/mvisonneau/automount/pulls).
