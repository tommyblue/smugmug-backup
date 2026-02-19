# Changelog

All notable changes to smugmug-backup will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## Unreleased

### Added

- Extended CSV metadata export with comprehensive data export in two separate files:
  - `albums_metadata.csv`: Complete album information including AlbumKey, Title, Description, Keywords, Privacy, SecurityType, ImageCount, dates, WebUri, and more (20 fields)
  - `images_metadata.csv`: Complete image/video information including filename, title, caption, format, dimensions, GPS coordinates, dates, and all technical metadata (31 fields)
- Added extensive fields to album and albumImage structures to capture all available SmugMug API data
- This provides a complete data export for users migrating away from SmugMug

### Changed

- Removed `HighlightImageUri` and `HighlightImageKey` fields from `albums_metadata.csv` as they contained incorrect data (highlight node IDs instead of actual image keys). The correct highlight image identification is available in the `IsHighlight` field of `images_metadata.csv`.

### Removed

-

### Fixed

- Fixed incorrect highlight/cover image identification in CSV export. The `IsHighlight` field in `images_metadata.csv` now correctly identifies the actual album cover image by properly parsing the SmugMug highlight node API response.

### Maintenance

-

## [v1.6.0](https://github.com/tommyblue/smugmug-backup/tree/v1.6.0) - 2024-12-13

### Added

- Add `FileNameNoExt`, `Extension`, `Date` and `Time` support in `store.file_names`

## [v1.5.5](https://github.com/tommyblue/smugmug-backup/tree/v1.5.5) - 2024-08-27

### Fixed

- skip videos that failed processing. They cannot be downloaded even from the website

### Maintenance

- Go v1.23
- Dependencies updates

## [v1.5.1](https://github.com/tommyblue/smugmug-backup/tree/v1.5.1) - 2023-09-15

### Fixed

- Do not replace some HTML characters in file name and other values used in the `store.file_names` conf

## [v1.5.0](https://github.com/tommyblue/smugmug-backup/tree/v1.5.0) - 2023-08-28

### Added

- Add concurrent file download and user album analysis with new `store.concurrent_downloads = <int>` and `store.concurrent_albums = <int>` configurations. These options can dramatically reduce backup speed with good connections. In my case a full backup (~200GB) dropped from ~12h to ~90m using 10 downloaders and 10 analyzers!! (both options defaults to 1 to maintain old behaviour).
- Add `store.force_video_download = <bool>` configuration to download videos marked as "Under process" in Smugmug. Sometimes those files will be permanently marked as such, so this option is the only way to backup them (or at least try to...)
- Add new `-cfg <folder path>` command line argument to add a custom folder containing a `config.toml` configuration file
- Print execution time to INFO log

## [v1.4.0](https://github.com/tommyblue/smugmug-backup/tree/v1.4.0) - 2023-08-06

### Added

- store files metadata to csv [#26]

### Maintenance

- Updated dependencies
- Update Dockerfile to ubuntu 22.04 and go v1.20.7
- Update goreleaser

## Fixed

- fix exception when album is empty [#39]

## [v1.3.1](https://github.com/tommyblue/smugmug-backup/tree/v1.3.1) - 2023-05-03

### Maintenance

- Update dependencies

## [v1.3.0](https://github.com/tommyblue/smugmug-backup/tree/v1.3.0) - 2023-02-27

### Added

- Add `-stats` flag to enable [statsviz](https://github.com/arl/statsviz) https://github.com/tommyblue/smugmug-backup/pull/32

### Maintenance

- Updated dependencies
- Update go to 1.20.1

## [v1.2.2](https://github.com/tommyblue/smugmug-backup/tree/v1.2.2) - 2020-11-27

### Fixed

- Fix `destination` path under Windows

## [v1.2.1](https://github.com/tommyblue/smugmug-backup/tree/v1.2.1) - 2020-09-25

### Maintenance

- Change how the image DateTimeOriginal field is parsed to manage wrong dates

## [v1.2.0](https://github.com/tommyblue/smugmug-backup/tree/v1.2.0) - 2020-09-24

### Added

- Add `store.use_metadata_times` and `store.force_metadata_times` confs to set files chtimes

### Changed

-

### Removed

- [DEPRECATION] `authentication.username` configuration is now ignored. The username is retrieved automatically from the API

### Maintenance

- Upgrade to Go 1.15.2

## [v1.1.1](https://github.com/tommyblue/smugmug-backup/tree/v1.1.1) - 2020-08-12

### Added

- Build releases for ARM and ARM64

### Maintenance

- Upgrade to Go 1.15

## [v1.1.0](https://github.com/tommyblue/smugmug-backup/tree/v1.1.0) - 2020-08-11

### Added

- `[store.file_names]` configuration added

## [v1.0.1](https://github.com/tommyblue/smugmug-backup/tree/v1.0.1) - 2020-08-10

### Maintenance

- This is an empty release used to create a version > v1.0.0 because an old temporary v1.0.0
  (that was removed) is the last one seen by listings that already pulled it before removal.
  That was my mistake but there's no way at the moment to retreat an already published version
  so this is the only solution (until go 1.16 that will introduce a way to retreat versions)

## [v0.0.4](https://github.com/tommyblue/smugmug-backup/tree/v0.0.4) - 2020-08-09

### Added

- Add configuration file support
- Add tests
- Add package docs: https://pkg.go.dev/github.com/tommyblue/smugmug-backup?tab=doc
- `-version` flag prints build version

### Changed

- **[BREAKING CHANGE]** Rename environment variables
- Move `main` package to `./cmd/smugmug-backup`

### Maintenance

- Update go version to 1.14.7
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
