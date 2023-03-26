# clq development
## Automated Builds
### Continuous Integration
Push to main, push to any other branches and pull-requests trigger a GitHub Action.
All these actions validate that the project conforms to our
quality standards: all tests pass, the docker images build and code coverage is sufficient.

A succesful push to main tags the commit and creates a GitHub release.

### Continuous Delivery
A succesful push to main builds and publishes the _latest_ docker image on Docker Hub.

The GitHub Action attaches binaries for unix, macos, and windows to the GitHub release
as well as builds and publishes docker images to dockerhub with the tagged version.

## Local Build
Run `make` to perform all the tests.
Run `make docker-test` to build and test the docker images.

## Architecture
![Class diagram](http://www.plantuml.com/plantuml/proxy?cache=no&src=https://raw.githubusercontent.com/denisa/clq/main/docs/class_diagram.puml)

## Processes
### Updating the version of golang
- `brew unpin go; brew upgrade go; brew pin go`
- `go.mod`
- `.github/workflows/ci.yaml`
- base images in `build/docker/alpine/Dockerfile`,  `build/docker/slim/Dockerfile`
Finally, run `make clean lcov`.

### Updating Dependencies
When merging a Dependabot change of depdency, run `make clean test lcov` to 
ensure that _all_ dependencies have been preserved. Dependabot tends to suppress
tool dependencies, for example `gcov2lcov`.
