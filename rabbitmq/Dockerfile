FROM alpine

ARG TARGETARCH
RUN CGO_ENABLED=0 GOOS=linux GOARCH=$TARGETARCH go build -mod=readonly ./pkg/...
COPY proxy-identity proxy-identity
RUN CGO_ENABLED=0 GOOS=linux GOARCH=$TARGETARCH go build -o /out/proxy-identity -mod=readonly -ldflags "-s -w" ./proxy-identity

COPY target/linux/mq_test .
COPY passwd.minimal /etc/passwd
RUN chown rabbit mq_test
USER rabbit
CMD ["sh", "-c", "tail -f /dev/null"]