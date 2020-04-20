# The base container contains all of the dependencies required to build, test,
# and run Cape.
FROM "golang:1.14.0-alpine"

ENV PATH="/usr/local/bin:${PATH}"
WORKDIR /go/src/github.com/capeprivacy/cape

# Install base dependencies in a different layer so they can be cached
# independent of the go related bits.
RUN apk --no-cache add make gcc musl-dev curl protoc protobuf-dev git docker

# Then we address our actual files that we're copying in - this way we only
# have to rebuild the system dep layer if we need to change the system
# dependencies
COPY go.mod go.sum Makefile tools.go ./
RUN make bootstrap

ENTRYPOINT ["sh", "-c"]