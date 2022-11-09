FROM rclone/rclone:1.60.0 AS rclone

FROM golang:1.19 AS builder

ENV CGO_ENABLED=0

WORKDIR /go/src/github.com/RoyXiang/putcallback/

COPY . .

RUN go install -v -ldflags "-s -w" -trimpath

FROM gcr.io/distroless/base-debian11:nonroot

LABEL \
	org.opencontainers.image.authors="developer@royxiang.me" \
	org.opencontainers.image.source="https://github.com/RoyXiang/putcallback"

COPY --from=builder --chown=nonroot /go/bin/putcallback /usr/local/bin/
COPY --from=rclone --chown=nonroot /usr/local/bin/rclone /usr/local/bin/

USER nonroot

EXPOSE 1880

VOLUME /home/nonroot/.config/rclone

ENTRYPOINT ["/usr/local/bin/putcallback"]
