# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [0ver](https://0ver.org).

## [Unreleased]

## [v0.1.2] - 2022-02-11

### Added

- `quay.io` releases
  
### Changed

- Bumped golang to `1.17`
- Bumped most dependencies

### Removed

- Removed brew/tap releases

## [v0.1.1] - 2021-05-28

### Changed

- Updated go fom `1.12` to `1.16`
- Migrated CI to GitHub actions
- Added new lint tests and fixed lint issues
- Refactored codebase more or less according to golang best practices
- Updated all dependencies to their latest versions
- Fixed some unhandled error cases
- Implemented goreleaser to generate artifacts
- Moved default branch to `main`

## [0.1.0] - 2019-06-06

### Added

- Added `validate` function that can list the availability of dependencies
- Added flag `--reuse-formatted-devices` to be able to reformat unused/unmounted devices
- Added flag `--use-lvm` in order to leverage LVM for the partioning of the devices
- Added flag `--use-all-devices` to use all available devices
- Updated `--devices` to `--device` and allowed to specify it multiple times
- Added support for /dev/xvd* devices

### Changed

- Refactored and cleaned up code

## [0.0.1] - 2019-05-23

### Added

- Working state of the app
- Formatting and mounting capabilities
- Automatic lookup of available devices
- Configurable device name to mount
- fstab management
- Makefile
- License
- Readme

[Unreleased]: https://github.com/mvisonneau/automount/compare/v0.1.2...HEAD
[v0.1.2]: https://github.com/mvisonneau/automount/tree/v0.1.2
[v0.1.1]: https://github.com/mvisonneau/automount/tree/v0.1.1
[0.1.0]: https://github.com/mvisonneau/automount/tree/0.1.0
[0.0.1]: https://github.com/mvisonneau/automount/tree/0.0.1
