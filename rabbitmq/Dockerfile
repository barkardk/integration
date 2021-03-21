FROM golang:1.16 as build
ARG GOARCH=amd64
ARG GOOS=linux
ARG CGO_ENABLED=0
ARG BUILD_TIME=$(date -u '+%F_%T')
ARG VERSION=""

WORKDIR /src
COPY go.mod /src/
COPY it /src/it
RUN go mod download -json >/dev/null && \
    GO111MODULE=on CGO_ENABLED=$CGO_ENABLED GOARCH=$GOARCH GOOS=$GOOS go build -o target/mq_test -tags $VERSION -ldflags "-s -w -X main.Version=$VERSION -X main.BuildTime=$BUILD_TIME" ./it

FROM alpine
COPY --from=build /src/target/mq_test .
COPY passwd.minimal /etc/passwd
RUN chown rabbit mq_test
USER rabbit
CMD ["sh", "-c", "tail -f /dev/null"]