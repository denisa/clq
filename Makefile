.PHONY: test install cov lcov

test:
	go test -v ./... -covermode=count -coverprofile=coverage.out -coverpkg=./...

cov: test
	go tool cover -html=coverage.out

lcov: test bin/gcov2lcov
	@bin/gcov2lcov -infile=coverage.out -outfile=coverage.lcov

install: test
	go install ./...

bin/gcov2lcov:
	@env GOBIN=$$PWD/bin GO111MODULE=on go install github.com/jandelgado/gcov2lcov

