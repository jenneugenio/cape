# The base container contains all of the dependencies required to build, test,
# and run Cape.
FROM "golang:1.14.0-alpine"

ENV PATH="/usr/local/bin:${PATH}"
WORKDIR /go/src/github.com/capeprivacy/cape

# Install base dependencies in a different layer so they can be cached
# independent of the go related bits.
RUN apk --no-cache add gcc musl-dev curl

# Then we address our actual files that we're copying in - this way we only
# have to rebuild the system dep layer if we need to change the system
# dependencies
COPY go.mod go.sum tools.go ./
RUN go mod download && go install github.com/magefile/mage

# Now that the deps have been installed (and cached in the previous layer) we
# can copy in what we need to build our magefile & cache it.
COPY bootstrap.go magefile.go ./
COPY mage ./mage
RUN mage -compile ./bin/mage-binary

# Now that we have our built mage binary, we can boostrap the system :)
RUN ./bin/mage-binary bootstrap

ENTRYPOINT ["sh", "-c"]
