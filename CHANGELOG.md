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

## [1.6.4] - 2022-05-15
### Fixed
- Bump docker/login-action from 1 to 2

## [1.6.3] - 2022-05-08
### Fixed
- rename github branch to `main`

## [1.6.2] - 2022-05-07
### Fixed
- When clq disagrees with the version change, it now shows the correct version
  and the change kind responsible for it.
- Bump actions/checkout from 3.0.1 to 3.0.2
- Bump github.com/yuin/goldmark from 1.4.11 to 1.4.12
- Bump github/codeql-action from 1 to 2
- Bump alpine from 3.15.0 to 3.15.4 in /build/docker/alpine

## [1.6.1] - 2022-04-23
### Fixed
- Bump github.com/yuin/goldmark from 1.4.4 to 1.4.11
- Bump github.com/stretchr/testify from 1.7.0 to 1.7.1
- Bump actions/setup-go from 2.1.4 to 3.0.0
- Bump actions/checkout from 2.4.0 to 3.0.1
- Bump golang from 1.17.4 to 1.18.1

## [1.6.0] - 2021-12-09
### Changed
- Upgrade to go 1.17
- Bump golang in /build/docker/alpine
- Bump golang in /build/docker/slim

## [1.5.0] - 2021-12-04
### Changed
- Upgrade to go 1.16
- Bump golang in /build/docker/alpine
- Bump golang in /build/docker/slim
- Bump github.com/yuin/goldmark from 1.3.1 to 1.4.4
- Stop using Docker Content Trust (the golang images aren’t signed)

### Fixed
- Bump actions/setup-go from 2.1.3 to 2.1.4
- Bump alpine from 3.12 to 3.15.0 in /build/docker/alpine, /build/docker/slim
- Bump actions/checkout from 2.3.4 to 2.4.0
- Bump github.com/stretchr/testify from 1.6.1 to 1.7.0

## [1.4.1] - 2021-11-28
### Fixed
- DockerHub is not available anymore

## [1.4.0] - 2021-01-01
### Changed
- can read the mapping from categories to change from a json file

## [1.3.5] - 2020-12-28
### Fixed
- Bump github.com/yuin/goldmark from 1.2.1 to 1.3.1
- Bumps golang base image from 1.15.0-alpine3.12 to 1.15.6-alpine3.12
- Bump actions/checkout from v2.3.2 to v2.3.4
- Bump actions/create-release from v1.1.3 to v1.1.4
- Bump actions/setup-go from v2.1.2 to v2.1.3
- Bump coverallsapp/github-action from v1.1.1 to v1.1.
- workflow converted from set-env to set-output

## [1.3.4] - 2020-08-17
### Fixed
- Bump golang to 1.15
- Update actions/checkout requirement to v2.3.2
- Bump actions/create-release from v1.1.2 to v1.1.3
- Bump actions/setup-go from v2.1.1 to v2.1.2

## [1.3.3] - 2020-08-04
### Fixed
- dependabot was not re-configured for the re-organized Dockerfiles

## [1.3.2] - 2020-08-04
### Fixed
- Bump github.com/yuin/goldmark from 1.2.0 to 1.2.1

## [1.3.1] - 2020-07-28
### Fixed
- Create codeql-analysis.yml
- Bump actions/setup-go from v2.1.0 to v2.1.1
- Bump github.com/yuin/goldmark from 1.1.32 to 1.2.0

## [1.3.0] - 2020-07-23
### Changed
- provide a new alpine-based docker image for use as a parent image by the clq-action

## [1.2.1] - 2020-07-21
### Fixed
- add class diagram
- use CHANGELOG.md to fill in GitHub release information

## [1.2.0] - 2020-07-19
### Changed
- option '-output md' will format the result of a query as markdown.
  There is an implied '-output json' if left unspecified.
- query can return multiple results

### Fixed
- Bump actions/checkout from v2.2.0 to v2.3.1
- Bump actions/create-release from v1.1.0 to v1.1.2
- Bump actions/setup-go from v2.0.3 to v2.1.0

## [1.1.2] - 2020-06-09
### Fixed
- pull-request fails when a tag exists for the CHANGELOG’s release

## [1.1.1] - 2020-06-08
### Fixed
- recognize release status _prereleased_ when creating github’s release
- include latest version in readme

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
