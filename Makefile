.PHONY: test
test:
	go test -v ./... -covermode=count -coverprofile=coverage.out -coverpkg=./...

.PHONY: cov
cov: test
	go tool cover -html=coverage.out

.PHONY: lcov
lcov: test bin/gcov2lcov
	@bin/gcov2lcov -infile=coverage.out -outfile=coverage.lcov

.PHONY: assertVersionDefined
assertVersionDefined:
	test -n "${VERSION}" -a "${VERSION}" != "latest"

VERSION ?= latest
DIST=dist/
LDFLAGS=-ldflags="-w -s -extldflags '-static' -X main.version=${VERSION}"

${DIST}:
	mkdir -p ${DIST}

.PHONY: build
build: ${DIST}
	go build -o ${DIST}clq .

AMD64=darwin linux windows
TARGET_AMD64:=$(addprefix build-,${AMD64})
.PHONY: build-all ${TARGET_AMD64}
build-all: ${TARGET_AMD64}

${TARGET_AMD64}:build-%:
	CGO_ENABLED=0 GOOS=$* GOARCH=amd64 go build -a ${LDFLAGS} -o ${DIST}clq-$*-amd64 .

.PHONY: install
install: test
	go install ./...

.PHONY: clean
clean:
	go clean -i ./...
	rm -fr *.out *.lcov ${DIST} bin/

DOCKER=alpine slim
TARGET_DOCKER_BUILD:=$(addprefix docker-build-,${DOCKER})
.PHONY: docker-build ${TARGET_DOCKER_BUILD}
docker-build: ${TARGET_DOCKER_BUILD}
	docker tag denisa/clq:slim denisa/clq:latest

${TARGET_DOCKER_BUILD}:docker-build-%:
	docker build --build-arg DOCKER_TAG=${VERSION} --file build/docker/$*/Dockerfile -t denisa/clq:$* .

TARGET_DOCKER_TEST:=$(addprefix docker-test-,${DOCKER})
.PHONY: docker-test ${TARGET_DOCKER_TEST}
docker-test: ${TARGET_DOCKER_TEST}

${TARGET_DOCKER_TEST}:docker-test-%:docker-build-%
	docker-compose --file build/docker/$*/Dockerfile.test.yml up

.PHONY: docker-push
docker-push:
	docker tag denisa/clq:slim denisa/clq:${VERSION}
	docker tag denisa/clq:alpine denisa/clq:${VERSION}-alpine
	docker push denisa/clq:latest
	docker push denisa/clq:${VERSION}
	docker push denisa/clq:${VERSION}-alpine

bin/gcov2lcov:
	env GOBIN=$$PWD/bin GO111MODULE=on go install github.com/jandelgado/gcov2lcov
