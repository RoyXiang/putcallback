FROM golang:1.19 AS builder

ENV CGO_ENABLED=0

WORKDIR /go/src/github.com/RoyXiang/putcallback/

COPY . .

RUN go install -v -ldflags "-s -w" -trimpath

FROM gcr.io/distroless/base-debian11:nonroot

COPY --from=builder --chown=nonroot /go/bin/putcallback /usr/local/bin/

USER nonroot

EXPOSE 1880

ENTRYPOINT ["/usr/local/bin/putcallback"]
