# syntax=docker/dockerfile:1

FROM golang:1.16-alpine as builder
COPY . /build/
WORKDIR /build
RUN env GOOS=linux GARCH=amd64 CGO_ENABLED=0 \
    go build -ldflags "-s -w" -o app main.go

FROM alpine:3.14.2
RUN apk add --update --no-cache ca-certificates curl

COPY --from=builder /build/app /app


CMD ["/app"]
