FROM golang:1.20.2-alpine3.17 AS builder

ARG DOCKER_TAG=latest

RUN apk update && apk add --no-cache git ca-certificates tzdata && update-ca-certificates

WORKDIR $GOPATH/src/denisa/clq/
COPY . .
RUN \
  go get -d -v && \
  CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build \
      -ldflags="-X main.version=${DOCKER_TAG} -w -s -extldflags '-static'" -a \
      -o /go/bin/clq .

FROM alpine:3.17.3
COPY --from=builder /go/bin/clq /usr/bin/clq
ENTRYPOINT ["/usr/bin/clq"]
