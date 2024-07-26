FROM rclone/rclone:1.67.0 AS rclone

FROM golang:1.21 AS builder

ENV CGO_ENABLED=0

WORKDIR /go/src/github.com/RoyXiang/putcallback/

COPY . .

RUN go install -v -ldflags "-s -w" -trimpath

FROM gcr.io/distroless/base-debian11:nonroot

COPY --from=builder --chown=nonroot /go/bin/putcallback /usr/local/bin/
COPY --from=rclone --chown=nonroot /usr/local/bin/rclone /usr/local/bin/

USER nonroot

EXPOSE 1880

VOLUME /home/nonroot/.config/rclone

ENTRYPOINT ["/usr/local/bin/putcallback"]
