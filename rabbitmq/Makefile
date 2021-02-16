NON_VENDOR_PACKAGES=$(shell go list ./... | grep -v "vendor")
TAG=$(shell git  rev-parse --short @)
DOCKER_REGISTRY=ghcr.io/barkardk
fmt:
	go fmt $(NON_VENDOR_PACKAGES)
vet:
	go vet $(NON_VENDOR_PACKAGES)

docker.restart: docker.stop docker.start
build: build.linux build.docker.and.push

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
	env GO111MODULE=on CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -o target/linux/mq_test -tags prod -mod=readonly -ldflags "-s -w " ./it
	chmod 755 target/linux/mq_test
build.docker.and.push:
	docker build -t ${TAG} .
	docker tag ${TAG} ${DOCKER_REGISTRY}/rabbitmq-client:${TAG}
	docker tag ${TAG} ${DOCKER_REGISTRY}/rabbitmq-client:latest
	docker push ${DOCKER_REGISTRY}/rabbitmq-client:${TAG}
	docker push ${DOCKER_REGISTRY}/rabbitmq-client:latest
build.darwin:
	mkdir -p target/darwin
	go build -o target/darwin/mq_test -tags dev it/main.go
	chmod 755 target/darwin/mq_test





