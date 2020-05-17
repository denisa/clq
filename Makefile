.PHONY: test
test:
	go test -v ./... -covermode=count -coverprofile=coverage.out -coverpkg=./...

.PHONY: cov
cov: test
	go tool cover -html=coverage.out

.PHONY: lcov
lcov: test bin/gcov2lcov
	@bin/gcov2lcov -infile=coverage.out -outfile=coverage.lcov

.PHONY: install
install: test
	go install ./...

.PHONY: docker-build
docker-build:
	@export DOCKER_CONTENT_TRUST=1 && docker build -f Dockerfile -t denisa/clq .

.PHONY: docker-test
docker-test:
	@docker-compose -f Dockerfile.test.yml up

bin/gcov2lcov:
	@env GOBIN=$$PWD/bin GO111MODULE=on go install github.com/jandelgado/gcov2lcov
