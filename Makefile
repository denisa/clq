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
	test -n "${VERSION}"

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

.PHONY: build-docker
build-docker:
	export DOCKER_CONTENT_TRUST=1 && docker build -f Dockerfile -t denisa/clq .

.PHONY: docker-test
docker-test:
	docker-compose -f Dockerfile.test.yml up

bin/gcov2lcov:
	env GOBIN=$$PWD/bin GO111MODULE=on go install github.com/jandelgado/gcov2lcov
