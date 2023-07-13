# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic
Versioning](http://semver.org/spec/v2.0.0.html).

## Unreleased

## [1.1.0] - 2023-07-14

### Added
- Add command line option `include-check-name`

## [1.0.0] - 2023-07-12

### Changed
- Rename labels used in event status (for `sensu_event_*` metrics)
    - `entity` => `sensu_entity_name`
    - `check` => `sensu_check_name`

### Security
- Update dependencies in `go.mod`

## [0.2.0] - 2022-09-01

### Added
- Add support event status including

## [0.1.0] - 2022-07-12

### Added
- Initial release
