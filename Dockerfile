FROM golang:alpine as builder

WORKDIR /code
COPY . /code
RUN go build .

FROM alpine:latest
LABEL org.opencontainers.image.source="https://github.com/xbmlz/sonar-dingtalk-plugin"

RUN apk add --no-cache ca-certificates tzdata
ENTRYPOINT ["/sonar-dingtalk-plugin"]