NON_VENDOR_PACKAGES=$(shell go list ./... | grep -v "vendor")

ifndef VERSION
	VERSION ?= $(shell (git describe --tags --dirty --match='v*' 2>/dev/null || echo v0.0.0) | cut -c2-)
endif
ifndef DOCKER_REGISTRY
	DOCKER_REGISTRY=ghcr.io/barkardk
endif
fmt:
	go fmt $(NON_VENDOR_PACKAGES)
vet:
	go vet $(NON_VENDOR_PACKAGES)

docker.restart: docker.stop docker.start
build: build.linux build.docker.and.push
release: fmt vet build.linux build.docker.and.push

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





