NON_VENDOR_PACKAGES=$(shell go list ./... | grep -v "vendor")

ifndef VERSION
	VERSION ?= $(shell (git describe --tags --dirty --match='v*' 2>/dev/null || echo v0.0.0) | cut -c2-)
endif
ifndef DOCKER_REGISTRY
	DOCKER_REGISTRY=ghcr.io/barkardk
endif

docker.restart: docker.stop docker.start
build: build.linux
release: build.linux

docker.start:
	docker-compose up -d --remove-orphans;
docker.stop:
	docker-compose stop;
docker.build:
	docker build -t rabbitmq-client
test:
	go test -tags=rabbitintegration  it/mq_test.gp -v -count=1;
build.linux:
	mkdir -p target/linux
	env GO111MODULE=on CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -o target/linux/mq_test -tags ${VERSION} -mod=readonly -ldflags "-s -w -X main.version=$(git describe --always --long --dirty)" ./it
	chmod 755 target/linux/mq_test
build.docker.and.push:
	docker build -t rabbitmq-client:${VERSION} .
	docker tag rabbitmq-client:${VERSION} ${DOCKER_REGISTRY}/rabbitmq-client:${VERSION}
	docker push ${DOCKER_REGISTRY}/rabbitmq-client:${VERSION}
build.darwin:
	mkdir -p target/darwin
	go build -o target/darwin/mq_test -tags dev  it/main.go
	chmod 755 target/darwin/mq_test

.PHONY: test
test:
	@echo "-> $@"
	go test -v -race ./...

.PHONY: fmt
fmt:
	@echo "-> $@"
	@gofmt -s -l ./ | grep -v vendor | tee /dev/stderr

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




