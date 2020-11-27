CURRENT_DIR = $(patsubst %/,%,$(dir $(abspath $(lastword $(MAKEFILE_LIST)))))

GOPATH       = $(shell go env GOPATH)
CGO_ENABLED  = 0
GOOS        ?= linux
GOARCH      ?= amd64
GOFLAGS      = CGO_ENABLED=$(CGO_ENABLED) GOOS=$(GOOS) GOARCH=$(GOARCH)

WERF_PATH   = $(shell multiwerf werf-path 1.1 rock-solid)
WERF_CONFIG = $(CURRENT_DIR)/scripts/werf/werf.yaml

# Download dependencies.
.PHONY: gomod
gomod:
	@echo "+@"
	@go mod download

# Check lint, code styling rules. e.g. pylint, phpcs, eslint, style (java) etc ...
.PHONY: style
style:
	@echo "+ $@"
	@golangci-lint run -v

# Format code. e.g Prettier (js), format (golang)
.PHONY: format
format:
	@echo "+ $@"
	@go fmt "$(CURRENT_DIR)/..."

# Shortcut to launch all the test tasks (unit, functional and integration).
.PHONY: test
test: test-unit
	@echo "+ $@"

# Launch unit tests. e.g. pytest, jest (js), phpunit, JUnit (java) etc ...
.PHONY: test-unit
test-unit:
	@echo "+ $@"
	@go test \
		-race \
		-v \
		-cover \
		-coverprofile \
		coverage.out \
		"$(CURRENT_DIR)/..."

# Build binary file
.PHONY: go-build
go-build:
	@echo "+ $@"
	@$(GOFLAGS) go build \
		-ldflags "-s -w" \
		-o $(CURRENT_DIR)/out/afi-file-downloader \
		$(CURRENT_DIR)/cmd/afi-file-downloader/main.go

# Clean out directory
.PHONY: clean
clean:
	@echo "+ $@"
	@rm -rf $(CURRENT_DIR)/out

# Build docker image using werf
.PHONY: werf-build
werf-build:
	@echo "+ $@"
	@$(WERF_PATH) build \
		--config $(WERF_CONFIG) \
		--stages-storage :local

# Publish image
.PHONY: werf-publish
werf-publish:
	@echo "+ $@"
	@$(WERF_PATH) publish \
		--config $(WERF_CONFIG) \
		--stages-storage :local \
		--images-repo $(DOCKER_REPO) \
		--tag-by-stages-signature
