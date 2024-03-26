.PHONY: test
test:
	go test -v ./... -covermode=count -coverprofile=coverage.out -coverpkg=./...

.PHONY: cov
cov: test
	go tool cover -html=coverage.out

.PHONY: lcov
lcov: test bin/gcov2lcov
	@bin/gcov2lcov -infile=coverage.out -outfile=coverage.lcov

.PHONY: superlinter
superlinter:
	docker run --rm \
		--platform linux/amd64 \
		-e RUN_LOCAL=true \
		--env-file ".github/super-linter.env" \
		-w /workspace -v "$$PWD":/workspace \
		ghcr.io/super-linter/super-linter:v5

.PHONY: golint
golint:
	docker run -t --rm \
		-w /workspace -v "$$PWD":/workspace \
		golangci/golangci-lint:v1.57.1 golangci-lint \
		run --config .github/linters/.golangci.yml --verbose --fast

.PHONY: plantuml
plantuml:
	docker run -t --rm \
		-w /workspace -v "$$PWD":/workspace \
		plantuml/plantuml:1.2024.2 \
		-v -o /workspace/_site docs/

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
TARGET_AMD64:=$(addprefix build-amd64-,${AMD64})
ARM64=darwin linux
TARGET_ARM64:=$(addprefix build-arm64-,${ARM64})
.PHONY: build-all ${TARGET_AMD64} ${TARGET_ARM64}
build-all: ${TARGET_AMD64} ${TARGET_ARM64}

${TARGET_AMD64}:build-amd64-%:
	CGO_ENABLED=0 GOOS=$* GOARCH=amd64 go build -a ${LDFLAGS} -o ${DIST}clq-$*-amd64 .

${TARGET_ARM64}:build-arm64-%:
	CGO_ENABLED=0 GOOS=$* GOARCH=arm64 go build -a ${LDFLAGS} -o ${DIST}clq-$*-arm64 .

.PHONY: install
install: test
	go install ./...

.PHONY: clean
clean:
	go clean -i ./...
	rm -fr *.out *.lcov ${DIST} bin/ _site/

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

bin/gcov2lcov:
	env GOBIN=$$PWD/bin GO111MODULE=on go install github.com/jandelgado/gcov2lcov
