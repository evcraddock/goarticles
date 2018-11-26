
# FROM golang:alpine as build-env
# ADD . /go/src/github.com/evcraddock/goarticles/
# RUN go build -o /go/bin/goarticles-api /go/src/github.com/evcraddock/goarticles/cmd/goarticles-api/goarticles-api.go
FROM alpine
RUN  apk update && \
     apk add libc6-compat && \
     apk add ca-certificates

WORKDIR /
EXPOSE 8000
COPY app/goarticles-api /goarticles-api

ENTRYPOINT ./goarticles-api >/var/log/app.log
