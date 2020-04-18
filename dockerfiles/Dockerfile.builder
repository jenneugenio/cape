# This container is used to build a version of Cape without needing to have any
# of the local dependencies installed.
#
# It enables a user to volume mount their local repository of Cape and go build
# cache to make it easy (and fast) to rebuild versions of Cape.
#
# This container sits on top of the base cape container which manages .
FROM capeprivacy/base:latest

VOLUME /go/src/github.com/capeprivacy/cape
VOLUME /root/.cache

ENTRYPOINT ["sh", "-c"]
