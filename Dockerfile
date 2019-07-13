FROM alpine:3.8

ADD sidecar-inject-server /sidecar-inject-server
ENTRYPOINT ["./sidecar-inject-server"]