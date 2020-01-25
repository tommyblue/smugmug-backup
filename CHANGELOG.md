# Changelog

All notable changes to smugmug-backup will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## Unreleased

### Added

- 

### Changed

- 

### Removed

- 

### Maintenance

- 

## [v1.0.0](https://github.com/tommyblue/smugmug-backup/tree/v1.0.0) - 2020-01-25

First official release

### Features

- Download images and videos from a SmugMug account
- On subsequent executions, only download new files or exising files (same path and name) if the size has changed
- Does not delete local files if deleted in SmugMug
- Skip videos "under processing"
- Can download files without a name, using their SmugMug ID
