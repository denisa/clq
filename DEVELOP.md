# clq development

## Automated Builds

### Continuous Integration

Push to main, push to any other branches and pull-requests trigger a GitHub Action.
All these actions validate that the project conforms to our
quality standards: all tests pass, the Docker images build and code coverage is sufficient.

A successful push to main tags the commit and creates a GitHub release.

### Continuous Delivery

A successful push to main builds and publishes the *latest* Docker image on Docker Hub.

The GitHub Action attaches binaries for Unix, macOS, and windows to the GitHub release
as well as builds and publishes Docker images to Docker Hub with the tagged version.

## Local Build

Run `make cov` to perform all the tests.
Run `make docker-test` to build and test the Docker images.
Run `make super-linter` to lint the complete project.

## Architecture

![Class diagram](https://denisa.github.io/clq/class_diagram.png)

## Processes

### Updating the version of Go

- `brew unpin go; brew upgrade go; brew pin go`
- `go.mod`
- base images in `build/docker/alpine/Dockerfile`, `build/docker/slim/Dockerfile`
Finally, run `make clean test docker-test`.

### Docker: updating Alpine

Dependabot will update only the base image for `denisa/clq:slim`.
Manually update the alpine version used by the *builder* for [alpine](build/docker/alpine/Dockerfile)
and [slim](build/docker/slim/Dockerfile).

### Updating Super Linter

Dependabot will update the version in [linter.yaml](.github/workflows/linter.yaml), ensure that the
`super-linter` target in [Makefile](Makefile) is correct.

Also runs `make super-linter` to see if new fixes are needed.
