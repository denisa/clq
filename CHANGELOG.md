# changelog



[![Keep a Changelog](https://img.shields.io/badge/Keep%20a%20Changelog-1.0.0-informational)](https://keepachangelog.com/en/1.0.0/)
[![Semantic Versioning](https://img.shields.io/badge/Sematic%20Versioning-2.0.0-informational)](https://semver.org/spec/v2.0.0.html)
![clq validated](https://img.shields.io/badge/clq-validated-success)

Keep the newest entry at top, format date according to ISO 8601: `YYYY-MM-DD`.

Categories:
- _major_ release trigger:
   - `Added` for new features.
   - `Removed` for now removed features.
- _minor_ release trigger:
   - `Changed` for changes in existing functionality.
   - `Deprecated` for soon-to-be removed features.
- _bug-fix_ release trigger:
   - `Fixed` for any bug fixes.
   - `Security` in case of vulnerabilities.
   
## [1.1.1] - 2020-06-08
### Fixed
- recognize release status _prereleased_ when creating github’s release

## [1.1.0] - 2020-06-08
### Changed
- a query for a release status can now return _prereleased_

## [1.0.2] - 2020-06-08
### Fixed
- now set the version in the released binaries
- bump github.com/stretchr/testify from 1.6.0 to 1.6.1
- bump github.com/yuin/goldmark from 1.1.31 to 1.1.32

## [1.0.1] - 2020-06-07
### Fixed
- the _tag_ job will not run unless all tests are green

## [1.0.0] - 2020-06-05
### Added
- basic validations
- basic queries
- docker image and (amd64) binaries for linux, macos and windows
- dependabot configuration
