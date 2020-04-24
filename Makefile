.PHONY: test install cov

test:
	go test -v ./... -covermode=count -coverprofile=coverage.out -coverpkg=./...

cov: test
		go tool cover -html=coverage.out

install: test
	go install ./...
