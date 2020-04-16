PREFIX?=
GOOS_OVERRIDE?=
SRC=$(shell find . -type f -name '*.go')

all: ci

.PHONY: all

# ###############################################
# Required Binaries & Versioning Checks
#
# The following section lists out all of the system dependencies (and their
# versions) that this Makefile depends upon on-top of standard tooling such as
# Git.
# ###############################################
DOCKER_REQUIRED_VERSION=18.
GO_REQUIRED_VERSION=1.14.
HELM_REQUIRED_VERSION=v3.0.
GSUTIL_REQUIRED_VERSION=4.47
PROTOC_REQUIRED_VERSION=3.11.4

UNAME_S := $(shell uname -s)
ifeq ($(UNAME_S),Linux)
	PROTOC_ZIP=protoc-$(PROTOC_REQUIRED_VERSION)-linux-x86_64.zip
endif
ifeq ($(UNAME_S),Darwin)
	PROTOC_ZIP=protoc-$(PROTOC_REQUIRED_VERSION)-osx-x86_64.zip
endif

CURRENT_GO_VERSION := $(shell go version 2> /dev/null)
gocheck:
ifndef CURRENT_GO_VERSION
ifeq (,$(findstring $(GO_REQUIRED_VERSION),$(CURRENT_GO_VERSION)))
ifeq (,$(BYPASS_GO_CHECK))
	$(error "Go Version $(GO_REQUIRED_VERSION) is required, found $(CURRENT_GO_VERSION)")
endif
endif
endif

CURRENT_DOCKER_VERSION := $(shell docker version 2> /dev/null)
dockercheck:
ifeq (,$(DOCKER_PATH))
ifeq (,$(shell command -v docker 2> /dev/null))
ifeq (,$(findstring $(DOCKER_REQUIRED_VERSION),$(CURRENT_DOCKER_VERSION)))
ifeq (,$(BYPASS_DOCKER_CHECK))
	$(error "Docker version $(DOCKER_REQUIRED_VERSION) is required, found $(CURRENT_DOCKER_VERSION)")
endif
endif
endif
endif

CURRENT_HELM_VERSION := $(shell helm version 2> /dev/null)
helmcheck:
ifeq (,$(shell command -v helm 2> /dev/null))
ifeq (,$(findstring $(HELM_REQUIRED_VERSION),$(CURRENT_HELM_VERSION)))
ifeq (,$(BYPASS_HELM_CHECK))
	$(error "Helm version $(HELM_REQUIRED_VERSION) is required, found $(CURRENT_HELM_VERSION).")
endif
endif
endif

CURRENT_GSUTIL_VERSION := $(shell gsutil version 2> /dev/null)
gsutilcheck:
ifeq (,$(shell command -v gsutil 2> /dev/null))
ifeq (,$(findstring $(GSUTIL_REQUIRED_VERSION),$(CURRENT_GSUTIL_VERSION)))
ifeq (,$(BYPASS_GSUTIL_CHECK))
	$(error "gsutil version $(GSUTIL_REQUIRED_VERSION) is required, found $(CURRENT_GSUTIL_VERSION)")
endif
endif
endif

CURRENT_PROTOC_VERSION := $(shell protoc --version 2> /dev/null)
protoccheck:
ifeq (,$(shell command -v protoc 2> /dev/null))
ifeq (,$(findstring $(PROTOC_REQUIRED_VERSION),$(PROTOC_REQUIRED_VERSION)))
ifeq (,$(BYPASS_PROTOC_CHECK))
	$(error "protoc version $(PROTOC_REQUIRED_VERSION) is required, found $(PROTOC_REQUIRED_VERSION)")
endif
endif
endif

.PHONY: gocheck dockercheck

# ###############################################
# Version Derivation
#
# Rules and variable definitions used to derive the current version of the
# source code. This information is also used for deriving the type of release
# to perform if `make push` is invoked.
# ###############################################
VERSION=$(shell [ -d .git ] && git describe --tags --abbrev=0 2> /dev/null | sed 's/^v//')
EXACT_TAG=$(shell [ -d .git ] && git describe --exact-match --tags HEAD 2> /dev/null | sed 's/^v//')
ifeq (,$(VERSION))
    VERSION=dev
endif
NOT_RC=$(shell git tag --points-at HEAD | grep -v -e -rc)

ifeq ($(EXACT_TAG),)
    PUSHTYPE=master
else
    ifeq ($(NOT_RC),)
	PUSHTYPE=release-candidate
    else
	PUSHTYPE=release
    endif
endif

# ###############################################
# Setup and Teardown
#
# This section of the makefile focuses on the rules needed to setup an
# environment to develop with this codebase and then subsequently start from
# scratch.
# ###############################################

bootstrap: download install-tools
	$(info All dependencies and tooling have been installed!)

download: gocheck
	go mod download

install-tools: gocheck install-protoc download tools.go
	cat tools.go | grep _ | awk -F'"' '{print $$2}' | xargs -tI % go install %

SUDO=$(shell which sudo)
install-protoc:
	# normal protoc doesn't work with alpine, assume its installed via apk
	if [ ! -f "/etc/alpine-release" ]; then \
		curl -OL https://github.com/protocolbuffers/protobuf/releases/download/v$(PROTOC_REQUIRED_VERSION)/$(PROTOC_ZIP); \
		$(SUDO) unzip -o $(PROTOC_ZIP) -d /usr/local bin/protoc; \
		$(SUDO) unzip -o $(PROTOC_ZIP) -d /usr/local 'include/*'; \
	fi; \

clean: gocheck
	go clean
	rm $(PREFIX)bin/cape

gogen: gocheck protoccheck gqlgen.yml coordinator/schema/*.graphql
	go generate ./...

.PHONY: bootstrap clean gogen

bootstrap-local-dev: bootstrap-helm

bootstrap-helm: helm-install helm-add-stable helm-update

helm-install:
ifeq (, $(shell which helm))
	curl -fsSL -o get_helm.sh https://raw.githubusercontent.com/helm/helm/master/scripts/get-helm-3
	chmod 700 get_helm.sh
	./get_helm.sh
endif

helm-add-stable: helmcheck
	helm repo add stable https://kubernetes-charts.storage.googleapis.com

helm-update: helmcheck
	helm repo update

.PHONY: bootstrap-local-dev

# ###############################################
# Testing, Building and Formatting
#
# This section of the makefile focuses on the rules needed to test and build
# code contained within this repository.
# ###############################################

lint: gocheck
	golangci-lint run

BUILD_DATE=$(date)
PKG=github.com/capeprivacy/cape
DATE=$(shell date)
GO_BUILD=go build -i -v -ldflags "-w -X '$(PKG)/version.Version=$(VERSION)' -X '$(PKG)/version.BuildDate=$(DATE)' -s"
$(PREFIX)bin/cape: gocheck $(SRC) gogen
	$(GOOS_OVERRIDE) $(GO_BUILD) -o $@ $(PKG)/cmd

build: $(PREFIX)bin/cape

test: lint integration

unit: gocheck
	go test -v ./...

CAPE_DB_URL?="postgres://postgres:dev@localhost:5432/postgres?sslmode=disable"
CAPE_DB_MIGRATIONS?="$(shell pwd)/migrations"
CAPE_DB_TEST_MIGRATIONS?="$(shell pwd)/database/dbtest/migrations"
CAPE_DB_SEED_MIGRATIONS?="$(shell pwd)/tools/seed"
integration: gocheck
	CAPE_DB_URL=$(CAPE_DB_URL) CAPE_DB_SEED_MIGRATIONS=$(CAPE_DB_SEED_MIGRATIONS) CAPE_DB_TEST_MIGRATIONS=$(CAPE_DB_TEST_MIGRATIONS) CAPE_DB_MIGRATIONS=$(CAPE_DB_MIGRATIONS) go test -v ./... -tags=integration

fmt: gocheck
	gofmt -s -l -w $(SRC)

tidy: gocheck
	go mod tidy
	if [ -n "$(shell git status --untracked-files=no --porcelain)" ]; then \
		echo "Make sure to run and commit changes from go mod tidy"; \
		exit 1; \
	fi; \


ci: tidy lint build test docker

.PHONY: lint build fmt test ci integration

# ###############################################
# Building Docker Image
#
# Builds a docker image for TF Encrypted that can be used to deploy and
# test.
# ###############################################
DOCKER_BUILD=docker build -t capeprivacy/$(1):$(2) -f $(3) .
docker: dockerfiles/Dockerfile.base dockerfiles/Dockerfile.coordinator dockerfiles/Dockerfile.connector dockercheck
	$(call DOCKER_BUILD,cape,latest,dockerfiles/Dockerfile.base)
	$(call DOCKER_BUILD,coordinator,latest,dockerfiles/Dockerfile.coordinator)
	$(call DOCKER_BUILD,connector,latest,dockerfiles/Dockerfile.connector)
	$(call DOCKER_BUILD,update,latest,dockerfiles/Dockerfile.update)
	$(call DOCKER_BUILD,customer_seed,latest,tools/seed/Dockerfile.customer)

.PHONY: docker

# ###############################################
# Releasing Docker Images
#
# Using the docker build infrastructure, this section is responsible for
# authenticating to docker hub and pushing built docker containers up with the
# appropriate tags.
# ###############################################
DOCKER_TAG=docker tag capeprivacy/$(1):$(2) docker.pkg.github.com/capeprivacy/cape/$(1):$(3)
DOCKER_PUSH=docker push docker.pkg.github.com/capeprivacy/cape/$(1):$(2)

docker-tag: dockercheck
	$(call DOCKER_TAG,cape,latest,$(VERSION))
	$(call DOCKER_TAG,coordinator,latest,$(VERSION))
	$(call DOCKER_TAG,connector,latest,$(VERSION))
	$(call DOCKER_TAG,update,latest,$(VERSION))
	$(call DOCKER_TAG,customer_seed,latest,$(VERSION))

docker-push-tag: dockercheck
	$(call DOCKER_PUSH,cape,$(VERSION))
	$(call DOCKER_PUSH,coordinator,$(VERSION))
	$(call DOCKER_PUSH,connector,$(VERSION))
	$(call DOCKER_PUSH,update,$(VERSION))
	$(call DOCKER_PUSH,customer_seed,$(VERSION))

docker-push-latest: dockercheck
	$(call DOCKER_PUSH,cape,latest)
	$(call DOCKER_PUSH,coordinator,latest)
	$(call DOCKER_PUSH,connector,latest)
	$(call DOCKER_PUSH,update,latest)
	$(call DOCKER_PUSH,customer_seed,latest)

.PHONY: docker-login docker-push-latest docker-push-tag docker-tag

# ###############################################
# Targets for pushing docker images
#
# The following are that are called dependent on the push type of the release.
# They define what actions occur depending no whether this is simply a build of
# master (or a branch), release candidate, or a full release.
# ###############################################

# For all builds on the master branch, build the container
docker-push-master: docker

# For all builds on the master branch, with an rc tag
docker-push-release-candidate: releasecheck docker-push-master docker-login docker-tag docker-push-tag

# For all builds on the master branch with a release tag
docker-push-release: docker-push-release-candidate docker-push-latest

# This command calls the right docker push rule based on the derived push type
docker-push: docker-push-$(PUSHTYPE)

.PHONY: docker-push docker-push-release docker-push-release-candidate docker-push-master

releasecheck:
ifneq (yes,$(RELEASE_CONFIRM))
	$(error "Set RELEASE_CONFIRM=yes to really build and push release artifacts")
endif

helm-version-check:
ifeq (,$(shell grep -e $(VERSION) charts/connector/Chart.yaml 2> /dev/null))
	$(error "Version specified in charts/connector/Chart.yaml does not match $(VERSION)")
endif
ifeq (,$(shell grep -e $(VERSION) charts/coordinator/Chart.yaml 2> /dev/null))
	$(error "Version specified in charts/coordinator/Chart.yaml does not match $(VERSION)")
endif

publish: releasecheck helmcheck gsutilcheck helm-version-check
	mkdir -p local-dir
	gsutil rsync -d gs://dropout-helm-repo local-dir/
	helm package charts/connector
	helm package charts/coordinator
	cp *.tgz local-dir
	helm repo index local-dir/ --url https://dropout-helm-repo.storage.googleapis.com
	gsutil rsync -d local-dir/ gs://dropout-helm-repo

.PHONY: publish
