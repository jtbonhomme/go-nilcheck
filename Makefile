IMAGES_TAG   = ${shell git describe --tags --match 'v[0-9]*\.[0-9]*\.[0-9]*' 2> /dev/null || echo 'latest'}
GIT_COMMIT   ?= $(shell git rev-parse HEAD)
GIT_TAG      ?= $(shell git tag --points-at HEAD)
DIST_TYPE    ?= snapshot
BRANCH       ?= $(shell git rev-parse --abbrev-ref HEAD)
REPO	 ?= $(shell echo $(JOB_NAME) | cut -d/ -f2)
DATE	 ?= $(shell date -u +%FT%T%z)

PROJECT_NAME := go-nilcheck
BINARY_NAME  := nilcheck
PKG_ORG      := github.com/jtbonhomme/$(PROJECT_NAME)
CMD	  := cmd/$(BINARY_NAME)
PKG	  := $(PKG_ORG)/$(CMD)
GO_FILES     := $(shell find . -name '*.go')

GO	   := go
GOLANGCILINT := golangci-lint
GORELEASER   := goreleaser
GOFMT	:= gofmt
GOIMPORTS    := goimports
OS	   := $(shell uname -s)
GOOS	 ?= $(shell echo $OS | tr '[:upper:]' '[:lower:]')
GOARCH       ?= amd64

BUILD_OPTS   = -ldflags "-X $(PKG_ORG)/internal/version.Tag=$(IMAGES_TAG) \
	-X $(PKG_ORG)/internal/version.GitCommit=$(GIT_COMMIT) \
	-X $(PKG_ORG)/internal/version.BuildTime=$(DATE)"

# HELP =================================================================================================================
# This will output the help for each task
# thanks to https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
help: ; $(info $(M) Display makefile targets…) @ ## Display this help screen
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
.PHONY: help

linter: ; $(info $(M) Lint go source code…) @ ### check by golangci linter.
	$(GOLANGCILINT) -v --deadline 100s --skip-dirs docs run ./...
.PHONY: linter

test: linter ; $(info $(M) Executing tests…)@ ### run tests.
	$(GO) test -v -coverprofile=cover.out -race ./...
.PHONY: test

cover: test ; $(info $(M) Testing with code coverage…)@ ## Measure the test coverage.
	gocov convert cover.out | gocov-xml > cover.xml
	gocov convert cover.out | gocov-html > cover.html
	open cover.html
.PHONY: cover

run: ; $(info $(M) Runing program…) @ ### run program.
	$(GO) run $(PKG)
.PHONY: run

fmt: ; $(info $(M) Formatting golang code…) @ ## Format go code.
	$(GOFMT) -w -s $(GO_FILES)
	$(GOIMPORTS) -w $(GO_FILES)
.PHONY: fmt

install: build ; $(info $(M) Installing $(BINARY_NAME)…) @ ## Install the binary as executable.
	$(shell cp $(BINARY_NAME) $(GOPATH)/bin/)
.PHONY: install

clean: ; $(info $(M) Cleaning project…) @ ## Build the main program.
	rm -f coverage.out cover.xml cover.html
	rm -f $(BINARY_NAME)
.PHONY: clean

build: ; $(info $(M) Building program executable…) @ ## Build the main program.
	$(GO) build -o $(BINARY_NAME) \
		$(BUILD_OPTS) \
		$(CMD)/main.go

deps: ; $(info $(M) Testing with code coverage…) @  test ## Measure the test coverage.
	which golangci-lint || (go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.50.1)
	which goreleaser || (go install github.com/goreleaser/goreleaser@latest)
	which gocov || (go get -u github.com/axw/gocov/gocov)
	which gocov-xml || (go get -u github.com/AlekSi/gocov-xml)
	which gocov-html || (go get -u github.com/matm/gocov-html)
.PHONY: deps

version: ; $(info $(M) Fetching version…) @ ## Build version.
ifneq ($(GIT_TAG),)
	$(eval VERSION := $(GIT_TAG))
	$(eval VERSION_FILE := $(GIT_TAG))
else
	$(eval VERSION := $(subst /,-,$(BRANCH)))
	$(eval VERSION_FILE := $(GIT_COMMIT)-SNAPSHOT)
endif
	@test -n "$(VERSION)"
	$(info Building $(VERSION)/$(VERSION_FILE) on sha1 $(GIT_COMMIT))

get_version: version ; $(info $(M) Building version…) @  ## Display version.
	$(info $(VERSION_FILE))

release: version ; $(info $(M) Releasing …) @  ## Release the program.
ifneq ($(GIT_TAG),)
	$(GORELEASER) release --parallelism 2 --rm-dist
else
	$(GORELEASER) release --snapshot --parallelism 2 --rm-dist
endif

snapshot: version ; $(info $(M) Releasing …) @  ## Release the program as a snapshot.
	$(GORELEASER) release --snapshot --parallelism 2 --rm-dist

goreleaser: version ; $(info $(M) Running goreleaser…) @ ## Run go releaser.
	$(GORELEASER) --parallelism 2 --rm-dist


