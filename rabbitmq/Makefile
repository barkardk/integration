.DEFAULT_GOAL := help
NON_VENDOR_PACKAGES=$(shell go list ./... | grep -v "vendor")
BUILD_TIME := $(shell date -u '+%F_%T')
ifndef VERSION
	VERSION ?= $(shell (git describe --tags --dirty --match='v*' 2>/dev/null || echo v0.0.0) | cut -c2-)
endif
ifndef DOCKER_REGISTRY
	DOCKER_REGISTRY=ghcr.io/barkardk
endif

.PHONY: help
help: ## Help
	@grep -E '^[a-zA-Z\\._-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[32m%-30s\033[0m %s\n", $$1, $$2}'

docker.restart: docker.stop docker.start ## restart docker-compose
build: fmt lint dep imp build.linux ## Run a full build with linting , depencency testing and import
release: publish ## Build linux release containers
release.binary: build.linux  ## Build binaries for various linux architectures

docker.start: ## Start docker-compose
	docker-compose up -d --remove-orphans;
docker.stop: ## Stop docker-compose
	docker-compose stop;
docker.build: ## Build rabbitmq-client docker container
	docker build -t rabbitmq-client

.PHONY: clean
clean: ## Clean compiled binaries
	@echo "-> $@"
	@rm -rf target

# -- build local multiarch binaries --
TARGETARCH := amd64 arm arm64
.PHONY: build.linux
build.linux: $(TARGETARCH) ## Build multiarch linux binaries

VERSION=$(shell (git describe --tags --dirty --match='v*' 2>/dev/null || echo v0.0.0) | cut -c2- )

.PHONY: $(TARGETARCH)
ARCH = ${word 1, $@}
$(TARGETARCH): ## Build multiarch linux binaries
	@echo "-> $@"
	env GO111MODULE=on CGO_ENABLED=0 GOARCH=$(ARCH) GOOS=linux go build -o target/linux-$(ARCH)/mq_test -tags ${VERSION} -ldflags "-s -w -X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME}" ./it
	@chmod 755 target/linux-$(ARCH)/mq_test

# -- build and push docker containers --

PLATFORMS ?= linux/amd64 linux/386 linux/arm linux/arm64 darwin/amd64
temp = $(subst /, ,$@)
GOOS = $(word 1, $(temp))
GOARCH = $(word 2, $(temp))

.PHONY: $(PLATFORMS) manifests publish
publish: $(PLATFORMS) ## Build docker multiarch containers and push to registry
$(PLATFORMS):
	@echo "-> $@"
	@echo "-- Build containers --"
	docker build -t ${DOCKER_REGISTRY}/rabbitmq-client:$(GOOS)-$(GOARCH)-${VERSION} -t ${DOCKER_REGISTRY}/rabbitmq-client:$(GOOS)-$(GOARCH)-latest --build-arg GOOS=$(GOOS) --build-arg GOARCH=$(GOARCH) --build-arg VERSION=${VERSION} .
	@echo "-- Push containers to registry --"
	docker push ${DOCKER_REGISTRY}/rabbitmq-client:$(GOOS)-$(GOARCH)-${VERSION}
	docker push ${DOCKER_REGISTRY}/rabbitmq-client:$(GOOS)-$(GOARCH)-latest
	@echo "-- Generate manifests for multiarch support --"
	docker manifest create --amend "${DOCKER_REGISTRY}/rabbitmq-client:${VERSION}" "${DOCKER_REGISTRY}/rabbitmq-client:$(GOOS)-$(GOARCH)-${VERSION}"
	docker manifest annotate "${DOCKER_REGISTRY}/rabbitmq-client:${VERSION}"  "${DOCKER_REGISTRY}/rabbitmq-client:$(GOOS)-$(GOARCH)-${VERSION}" --os=$(GOOS) --arch=$(GOARCH)
	docker manifest create --amend "${DOCKER_REGISTRY}/rabbitmq-client:latest" "${DOCKER_REGISTRY}/rabbitmq-client:$(GOOS)-$(GOARCH)-latest"
	docker manifest annotate "${DOCKER_REGISTRY}/rabbitmq-client:latest"  "${DOCKER_REGISTRY}/rabbitmq-client:$(GOOS)-$(GOARCH)-latest" --os=$(GOOS) --arch=$(GOARCH)

manifests: ## Publish manifests for docker multiarch support
	@echo "-> $@"
	docker manifest push "${DOCKER_REGISTRY}/rabbitmq-client:${VERSION}"
	docker manifest push "${DOCKER_REGISTRY}/rabbitmq-client:latest"

# -- build development binary --
.PHONY: build.darwin
build.darwin:
	@echo "-> $@"
	mkdir -p target/darwin
	go build -o target/darwin/mq_test -tags dev  it/main.go
	chmod 755 target/darwin/mq_test

.PHONY: test
test:
	@echo "-> $@"
	go test -tags=rabbitintegration  it/mq_test.gp -v -count=1;

.PHONY: fmt
fmt:
	@echo "-> $@"
	gofmt -s -l ./ | grep -v vendor | tee /dev/stderr

.PHONY: lint
lint:
	@echo "-> $@"
	@go get -u golang.org/x/lint/golint
	@golint ./... | tee /dev/stderr
	@go get -u golang.org/x/tools/go/analysis/cmd/vet
	@go vet --all

.PHONY: dep
dep:
	@echo "-> $@"
	@go mod tidy
	@go get -u github.com/golang/dep/cmd/dep
	@dep init && dep ensure -vendor-only

.PHONY: imp
imp:
	@echo "-> $@"
	@goimport




