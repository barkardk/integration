FROM alpine
COPY target/linux/it.test .
COPY passwd.minimal /etc/passwd
RUN chown rabbit it.test

FROM scratch
COPY --from=0 it.test .
COPY passwd.minimal /etc/passwd
USER rabbit
ENTRYPOINT ["./it.test", "-test.v", "-test.timeout=30s", "-test.count=2"]