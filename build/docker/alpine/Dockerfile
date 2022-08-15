FROM golang:1.18.2-alpine3.16 AS builder

ARG DOCKER_TAG=latest

RUN apk update && apk add --no-cache git ca-certificates tzdata && update-ca-certificates

WORKDIR $GOPATH/src/denisa/clq/
COPY . .
RUN go get -d -v
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
      -ldflags="-X main.version=${DOCKER_TAG} -w -s -extldflags '-static'" -a \
      -o /go/bin/clq .

FROM alpine:3.16.2
COPY --from=builder /go/bin/clq /usr/bin/clq
ENTRYPOINT ["/usr/bin/clq"]
