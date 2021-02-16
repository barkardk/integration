FROM alpine
COPY target/linux/mq_test .
COPY passwd.minimal /etc/passwd
RUN chown rabbit mq_test
USER rabbit
CMD ["sh", "-c", "tail -f /dev/null"]