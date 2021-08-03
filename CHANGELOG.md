# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Changed
- Login is more user friendly. Errors not printed directly.

### Removed
- PSL Report: it is now in the Heimat Web UI

## [0.1.4] - 2021-04-26

### Added

- `stats` command gathers time spent on different categories: projects, tasks and desc  
  > See more: [Stats](./README.md#stats)
- `version` (CLI) version flag to get the installed version

### Changed

- all displayed time formats now adhere to the ISO format

## [0.1.3] - 2021-03-09

### Added

- Time Report accepts a PSL filter; the filter can be a substring of a name and it is case-insensitive
- `time show day` shows breaks between tracked times 

### Changed

- Time Report groups by PSL

## [0.1.2] - 2021-03-03

### Added

- Time Delete, works similar to Time Add (`time delete`)

### Fixed

- Check Login Status for CLI Commands

## [0.1.1] - 2021-03-02

### Added
- Time Report Generation CLI Options for HR (`heimat -report=times`)

### Changed
- Wording of Holidays to Vacation days in user profile (`profile`)

### Fix
- User could "escape" into home state without login


## [0.1.0] - 2021-02-31

### Added
- Login into Heimat
- Logout
- Show basic user profile
- Time Tracking with default, relative and absolute date (`time add`)
- Time report of a day with default, relative and absolute date (`time show day`)
- Time report of the current month (`time show month`)
- Time copy of given day to another with default, relative and absolute date (`time copy`)

