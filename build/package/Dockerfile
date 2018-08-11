
FROM golang:alpine as build-env
ADD ../.. /go/src/github.com/evcraddock/goarticles/
RUN go build -o /go/bin/goarticles-api /go/src/github.com/evcraddock/goarticles/cmd/goarticles-api/goarticles-api.go
FROM alpine
RUN  apk update && \
     apk add libc6-compat && \
     apk add ca-certificates
WORKDIR /app
COPY --from=build-env /go/bin /app/
COPY ../../gcp.json /app/gcp.json
ENTRYPOINT ./goarticles-api
EXPOSE 8080