FROM golang:alpine as builder

WORKDIR /sonar-dingtalk-plugin-src
COPY --from=tonistiigi/xx:golang / /
COPY . /sonar-dingtalk-plugin-src
RUN go mod download && \
    make docker && \
    mv ./bin/sonar-dingtalk-plugin /sonar-dingtalk-plugin

FROM alpine:latest
LABEL org.opencontainers.image.source="https://github.com/xbmlz/sonar-dingtalk-plugin"

RUN apk add --no-cache ca-certificates tzdata
COPY --from=builder /sonar-dingtalk-plugin /
EXPOSE 9010
ENTRYPOINT ["/sonar-dingtalk-plugin"]