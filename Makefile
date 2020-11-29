CURRENT_DIR = $(patsubst %/,%,$(dir $(abspath $(lastword $(MAKEFILE_LIST)))))

GOPATH      = $(shell go env GOPATH)
CGO_ENABLED = 0
GOOS        = linux
GOARCH      = amd64
GOFLAGS     = CGO_ENABLED=$(CGO_ENABLED) GOOS=$(GOOS) GOARCH=$(GOARCH)

DOCKER_FILE       = $(CURRENT_DIR)/scripts/docker/Dockerfile
DOCKER_REPOSITORY = docker.io
DOCKER_IMAGE      = afi-file-downloader
DOCKER_TAG        = latest
DOCKER_IMAGE_NAME = $(DOCKER_REPOSITORY)/$(DOCKER_IMAGE):$(DOCKER_TAG)

WERF_PATH           = $(shell multiwerf werf-path 1.1 rock-solid)
WERF_CONFIG         = $(CURRENT_DIR)/scripts/werf/werf.yaml
WERF_STAGES_STORAGE = :local
WERF_DOCKER_OPTIONS = "-i"

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

# Run binary
.PHONY: go-run
go-run:
	@echo "+ $@"
	@$(CURRENT_DIR)/out/afi-file-downloader

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
		--stages-storage $(WERF_STAGES_STORAGE)

# Run docker image using werf
.PHONY: werf-run
werf-run:
	@echo "+ $@"
	@$(WERF_PATH) run \
		--config $(WERF_CONFIG) \
		--stages-storage $(WERF_STAGES_STORAGE) \
		--docker-options $(WERF_DOCKER_OPTIONS)

# Publish docker image using werf
.PHONY: werf-publish
werf-publish:
	@echo "+ $@"
	@$(WERF_PATH) publish \
		--config $(WERF_CONFIG) \
		--stages-storage $(WERF_STAGES_STORAGE) \
		--images-repo $(DOCKER_REPOSITORY) \
		--tag-by-stages-signature

# Build docker image
.PHONY: docker-build
docker-build:
	@echo "+ $@"
	@docker build \
		--rm \
		-t $(DOCKER_IMAGE):$(DOCKER_TAG) \
		-f $(DOCKER_FILE) \
		.
	@docker tag $(DOCKER_IMAGE):$(DOCKER_TAG) $(DOCKER_IMAGE_NAME)

# Run docker image
.PHONY: docker-run
docker-run:
	@echo "+ $@"
	@docker run \
		--rm \
		-i \
		$(DOCKER_IMAGE_NAME) \
		afi-file-downloader

# Publish docker image
.PHONY: docker-publish
docker-publish:
	@echo "+ $@"
	@docker login
	@docker push \
		$(DOCKER_IMAGE_NAME)
