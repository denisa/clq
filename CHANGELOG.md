# changelog

[![Keep a Changelog](https://img.shields.io/badge/Keep%20a%20Changelog-1.0.0-informational)](https://keepachangelog.com/en/1.0.0/)
[![Semantic Versioning](https://img.shields.io/badge/Sematic%20Versioning-2.0.0-informational)](https://semver.org/spec/v2.0.0.html)
![clq validated](https://img.shields.io/badge/clq-validated-success)

Keep the newest entry at top, format date according to ISO 8601: `YYYY-MM-DD`.

Categories, defined in [changemap.json](.github/clq/changemap.json)):

- *major* release trigger:
  - `Changed` for changes in existing functionality.
  - `Removed` for now removed features.
- *minor* release trigger:
  - `Added` for new features.
  - `Deprecated` for soon-to-be removed features.
- *bugfix* release trigger:
  - `Fixed` for any bugfixes.
  - `Security` in case of vulnerabilities.

## [1.8.7] - 2024-03-26

### Fixed

- Enable golint in module mode and provide the Makefile target `golint`
- Created [399](https://github.com/denisa/clq/issues/399) to investigate the suppressed cyclomatic complexity warning.
- Updated plantuml to 1.2024.2 and provide the Makefile target `plantuml`
- All generated files into `out/`

## [1.8.6] - 2024-03-12

### Fixed

- Bump peter-evans/dockerhub-description from 3 to 4
- Bump super-linter/super-linter from 5 to 6 while disabling go_modules and checkov;
  Created [394](https://github.com/denisa/clq/issues/394), [395](https://github.com/denisa/clq/issues/395)
  to address this later.
- Bump github.com/yuin/goldmark from 1.6.0 to 1.7.0
- Bump github.com/stretchr/testify from 1.8.4 to 1.9.0
- Requires golang 1.22.0
- Bumps alpine from 3.18.5 to 3.19.1

## [1.8.5] - 2023-12-30

### Fixed

- Bump github/codeql-action from 2 to 3
- Bump actions/checkout from 3 to 4 — the plantuml action was still on an old version
- Bump actions/upload-pages-artifact from 1 to 3
- Bump actions/deploy-pages from 1 to 4

## [1.8.4] - 2023-12-28

### Fixed

- Clarify handling of 'build' change kinds (out of spec!)

## [1.8.3] - 2023-12-15

### Fixed

- Error strings should not be capitalized or end with punctuation.
- Ensure consistent receiver type for Changelog and ChangeKind

## [1.8.2] - 2023-12-12

### Fixed

- Renders puml to GitHub page
- Resolve more warnings raised by Fleet
- Bumps actions/setup-go from 4 to 5

## [1.8.1] - 2023-12-04

### Fixed

- Bumps golang from 1.20 to 1.21
- Bumps alpine from 3.18.4 to 3.18.5
- Resolve various warnings raised by Fleet

## [1.8.0] - 2023-10-17

### Added

- It is now possible to define emoji to be included in the result. Makes for nicer release notes.

### Fixed

- `ioutil.ReadAll`, `ioutil.ReadFile` deprecated since golanmg 1.16
- Fix variables (config, require) that comflict with a package import
- Bump github.com/yuin/goldmark from 1.5.6 to 1.6.0
- Bumps alpine from 3.18.2 to 3.18.4.

## [1.7.16] - 2023-10-17

### Fixed

- As per [GitHub Blog](https://github.blog/changelog/2021-09-27-showing-code-scanning-alerts-on-pull-requests/),
  perform CodeQL analyses only for pushes to the `main` branch and pull requests agains the `main`.
- As
  per [GitHub Docs](https://docs.github.com/en/code-security/code-scanning/creating-an-advanced-setup-for-code-scanning/customizing-your-advanced-setup-for-code-scanning#avoiding-unnecessary-scans-of-pull-requests),
  skip CodeQL analyses for text files.

## [1.7.15] - 2023-10-16

### Fixed

- Document usage with [Whalebrew](https://github.com/whalebrew/whalebrew).

## [1.7.14] - 2023-10-07

### Fixed

- Use peter-evans/dockerhub-description to keep the description in DockerHub up-to-date

## [1.7.13] - 2023-10-07

### Fixed

- Update from github/super-linter@v5 to super-linter/super-linter@v5

## [1.7.12] - 2023-10-07

### Fixed

- Bump docker/login-action from 2 to 3
- Bump docker/setup-buildx-action from 2 to 3
- Bump docker/setup-qemu-action from 2 to 3
- Bump docker/build-push-action from 4 to 5
- Bump github.com/yuin/goldmark from 1.5.4 to 1.5.6

## [1.7.11] - 2023-09-04

### Fixed

- Bumps github/super-linter from 4 to 5
- Bumps actions/checkout from 3 to 4
- Bump github.com/stretchr/testify from 1.8.2 to 1.8.4
- Bump alpine from 3.17.2 to 3.18.2
- Bump golang from 1.20.2 to 1.20.5

## [1.7.10] - 2023-04-03

### Fixed

- clq now parses and renders links in the changelog

## [1.7.9] - 2023-04-02

### Fixed

- Provides better error message when the date is not correctly formatted
- Multi-line output with HERE doc (<https://github.com/github/docs/issues/21529>)
- GitHub Action runner advises against displaying the environment’s URL  
  `Warning: Skip setting environment url as environment 'dockerhub' may contain secret.`

## [1.7.8] - 2023-03-26

### Fixed

- Introduce super-linter

## [1.7.7] - 2023-03-25

### Fixed

- Add go report card, fix most issues

## [1.7.6] - 2023-03-25

### Fixed

- Bump golang from 1.19 to 1.20

## [1.7.5] - 2023-03-19

### Fixed

- Bump actions/setup-go from 3 to 4
- Bump docker/build-push-action from 3 to 4
- Bump coverallsapp/github-action from 1.1.3 to 2
- Bump github.com/stretchr/testify from 1.8.1 to 1.8.2
- Bump github.com/yuin/goldmark from 1.5.3 to 1.5.4
- Bump alpine from 3.17.0 to 3.17.2 in /build/docker/alpine
- Always checkout the repository first, this allows setup actions to uses info from the  
  repository. Looking at you, setup-go

## [1.7.4] - 2023-03-19

### Fixed

- Ask dependabot to only check for major version changes in GitHub actions.
- Properly handles soft and hard line breaks

## [1.7.3] - 2022-12-25

### Fixed

- Better error message when the changelog omits the Introduction’s title.
- Uses ncipollo/release-action v1, do not specify complete version
- Bump alpine from 3.16.2 to 3.17.0
- Bump github.com/stretchr/testify from 1.8.0 to 1.8.1
- Bump github.com/yuin/goldmark from 1.4.13 to 1.5.3
- Bump golang from 1.19.0-alpine3.16 to 1.19.4-alpine3.17

## [1.7.2] - 2022-09-05

### Fixed

- Dockerfile had hard-coded the target architecture, rely instead on Docker’s
  own `TARGETOS` and `TARGETARCH`

## [1.7.1] - 2022-08-30

### Fixed

- Various typo and inconsistencies in the documentation.

## [1.7.0] - 2022-08-29

### Added

- The docker images now support both `linux/amd64` and `linux/arm64`.
- Uses `docker/build-push-action`, remove docker-push target from the Makefile

### Fixed

- The make target `docker-test` now fails if the tested image produce an error
- The version of clq build and used by CI to extract the version information was being
  erroneously published to the release; stop doing that.

## [1.6.6] - 2022-08-27

### Fixed

- Release job was failing to upload artifacts because workflows had both
  ncipollo/release-action and actions/create-release
- Retire actions/upload-release-asset and let ncipollo/release-action upload artifacts
- Produce arm64 binaries for darwin and linux
- Bumps actions/setup-go from 3.2.1 to 3.3.0

## [1.6.5] - 2022-08-27

### Fixed

- Bump github.com/stretchr/testify from 1.7.1 to 1.8.0
- Bump github.com/yuin/goldmark from 1.4.12 to 1.4.13
- Bump actions/setup-go from 3.1.0 to 3.2.1
- Bump alpine from 3.15.4 to 3.16.2 in /build/docker/alpine
- Bump golang from 1.18.2 to 1.19.0

## [1.6.4] - 2022-05-15

### Fixed

- Uses denisa/semantic-tag-helper@v1
- Uses ncipollo/release-action
- Bump docker/login-action from 1 to 2
- Bump docker/setup-buildx-action from 1 to 2
- Bumps golang from 1.18.1-alpine3.15 to 1.18.2-alpine3.15.
- Bump actions/setup-go from 3.0.0 to 3.1.0
- Use a custom changemap in which Changes implies change to existing functionality, hence
  a major version bump.

## [1.6.3] - 2022-05-08

### Fixed

- rename GitHub branch to `main`

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

### Added

- Upgrade to go 1.17
- Bump golang in /build/docker/alpine
- Bump golang in /build/docker/slim

## [1.5.0] - 2021-12-04

### Added

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

### Added

- Can read the mapping from categories to change from a json file

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

- Dependabot was not re-configured for the re-organized Dockerfiles

## [1.3.2] - 2020-08-04

### Fixed

- Bump github.com/yuin/goldmark from 1.2.0 to 1.2.1

## [1.3.1] - 2020-07-28

### Fixed

- Create codeql-analysis.yml
- Bump actions/setup-go from v2.1.0 to v2.1.1
- Bump github.com/yuin/goldmark from 1.1.32 to 1.2.0

## [1.3.0] - 2020-07-23

### Added

- Provide a new alpine-based docker image for use as a parent image by the clq-action

## [1.2.1] - 2020-07-21

### Fixed

- Add class diagram
- Use CHANGELOG.md to fill in GitHub release information

## [1.2.0] - 2020-07-19

### Added

- Option '-output md' will format the result of a query as Markdown.
  There is an implied '-output json' if left unspecified.
- Query can return multiple results

### Fixed

- Bump actions/checkout from v2.2.0 to v2.3.1
- Bump actions/create-release from v1.1.0 to v1.1.2
- Bump actions/setup-go from v2.0.3 to v2.1.0

## [1.1.2] - 2020-06-09

### Fixed

- Pull-request fails when a tag exists for the CHANGELOG’s release

## [1.1.1] - 2020-06-08

### Fixed

- Recognize release status *prereleased* when creating github’s release
- Include latest version in readme

## [1.1.0] - 2020-06-08

### Added

- A query for a release status can now return *prereleased*

## [1.0.2] - 2020-06-08

### Fixed

- Now set the version in the released binaries
- Bump github.com/stretchr/testify from 1.6.0 to 1.6.1
- Bump github.com/yuin/goldmark from 1.1.31 to 1.1.32

## [1.0.1] - 2020-06-07

### Fixed

- the *tag* job will not run unless all tests are green

## [1.0.0] - 2020-06-05

### Added

- Basic validations
- Basic queries
- Docker image and (amd64) binaries for linux, macOS and windows
- Dependabot configuration
