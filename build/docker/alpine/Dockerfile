FROM dhi.io/golang:1.26.1-alpine3.23-dev AS builder
ARG DOCKER_TAG=latest

RUN apk add --no-cache git

WORKDIR $GOPATH/src/denisa/clq/
COPY . .
RUN \
  CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build \
      -ldflags="-X main.version=${DOCKER_TAG} -w -s -extldflags '-static'" -a \
      -o /go/bin/clq .

FROM dhi.io/alpine-base:3.23.3
COPY --from=builder /go/bin/clq /usr/bin/clq
ENTRYPOINT ["/usr/bin/clq"]
