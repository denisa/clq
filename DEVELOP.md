# clq development
## Automated Builds
### Continuous Integration
Push to master, push to any other branches and pull-requests trigger a GitHub Action.
All these actions validate that the project conforms to our
quality standards: all tests pass and code coverage is sufficient.

Pull-requests also trigger a test build on Docker Hub. This builds and tests the docker image.

A succesful push to master tags the commit and creates a GitHub release.

### Continuous Delivery
A succesful push to master builds and publishes the _latest_ docker image on Docker Hub.

The GitHub Action attaches binaries for unix, macos, and windows to the GitHub release
as well as builds and publishes a docker image on dockerhub with the tagged version.

## Local Build
Run `make` to perform all the tests.

## Architecture
![Class diagram](http://www.plantuml.com/plantuml/proxy?cache=no&src=https://raw.githubusercontent.com/denisa/clq/master/docs/class_diagram.puml)
