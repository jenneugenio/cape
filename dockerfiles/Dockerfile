FROM golang:1.13.7-alpine

WORKDIR /go/src/github.com/dropoutlabs/privacyai

COPY go.mod go.sum Makefile tools.go ./
RUN apk --no-cache add make gcc musl-dev \
    && make bootstrap

COPY . ./
RUN make build
