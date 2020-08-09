# Changelog

All notable changes to smugmug-backup will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## Unreleased

### Added

- Add configuration file
- Add tests
- Add package docs: https://pkg.go.dev/github.com/tommyblue/smugmug-backup?tab=doc

### Changed

- *BREAKING CHANGE* Rename environment variables
- Move `main` package to `./cmd/smugmug-backup`
- 

### Removed

- 

### Maintenance

- Update Dockerfile to go 1.14.7
- General refactoring of the code to make it more testable
- Improve documentation

## [v0.0.3](https://github.com/tommyblue/smugmug-backup/tree/v0.0.3) - 2020-07-25

### Added

- Windows support. Also build Windows binary on [release](https://github.com/tommyblue/smugmug-backup/releases)

### Maintenance

- Add the album path to the images, so it's possible to better debug images errors

## [v0.0.2](https://github.com/tommyblue/smugmug-backup/tree/v0.0.2) - 2020-07-01

### Changed

- almost all functions now return an error instead of `panic` or `log.Fatal`
- `ignorefetcherrors` command line option has been removed as now it's the default behaviour

### Maintenance

- Multiple refactorings

## [v0.0.1](https://github.com/tommyblue/smugmug-backup/tree/v0.0.1) - 2020-06-26

First release

### Features

- Download images and videos from a SmugMug account
- On subsequent executions, only download new files or exising files (same path and name) if the size has changed
- Does not delete local files if deleted in SmugMug
- Skip videos "under processing"
- Can download files without a name, using their SmugMug ID
