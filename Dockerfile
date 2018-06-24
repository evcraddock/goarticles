
FROM golang:alpine as build-env
ADD . /go/src/github.com/evcraddock/goarticles/
RUN cd /go/src/github.com/evcraddock/goarticles && go build -o goarticles
FROM alpine
RUN  apk update && \
     apk add libc6-compat && \
     apk add ca-certificates
WORKDIR /app
COPY --from=build-env /go/src/github.com/evcraddock/goarticles /app/
ENTRYPOINT ./goarticles
EXPOSE 8080