VERSION ?= latest
DIST=${OUT}dist/
OUT=out/
SITE=${OUT}site
PLANTUML_VERSION := $(shell grep -e '- uses: docker://plantuml/plantuml:' .github/workflows/plantuml.yaml | cut -d':' -f4)
SUPER_LINTER_VERSION := $(shell grep -e '- uses: super-linter/super-linter@' .github/workflows/linter.yaml | cut -d'@' -f2)

${DIST}:
	mkdir -p ${DIST}

${OUT}:
	mkdir -p ${OUT}

${SITE}:
	mkdir -p ${SITE}

.PHONY: clean
clean:
	go clean -i ./...
	rm -fr ${OUT}

.PHONY: test
test: ${OUT}
	go test -v ./... -covermode=count -coverprofile=${OUT}coverage.out -coverpkg=./...

.PHONY: cov
cov: test
	go tool cover -html=${OUT}coverage.out

.PHONY: super-linter
super-linter:
	docker run \
		--platform linux/amd64 \
		--rm \
		-e RUN_LOCAL=true \
		-e SHELL=/bin/bash \
		--env-file ".github/super-linter.env" \
		-v "$$PWD":/tmp/lint \
		ghcr.io/super-linter/super-linter:$(SUPER_LINTER_VERSION)

.PHONY: golint
golint:
	docker run -t --rm \
		-w /workspace -v "$$PWD":/workspace \
		golangci/golangci-lint:v1.57.1 golangci-lint \
		run --config .github/linters/.golangci.yml --verbose --fast

.PHONY: plantuml
plantuml: ${SITE}
	docker run -t --rm \
		-w /workspace -v "$$PWD":/workspace \
		plantuml/plantuml:$(PLANTUML_VERSION) \
		-v -o /workspace/${SITE} docs/

.PHONY: assertVersionDefined
assertVersionDefined:
	test -n "${VERSION}" -a "${VERSION}" != "latest"

LDFLAGS=-ldflags="-w -s -extldflags '-static' -X main.version=${VERSION}"

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
	LABEL=$* docker compose --file build/docker/compose.yaml up --exit-code-from sut
