.PHONY: build
build:
	GOOS=linux GOARCH=amd64 go build -o .build/goarticles-api cmd/goarticles-api/goarticles-api.go
	GOOS=linux GOARCH=amd64 go build -o ${GOPATH}/bin/goarticles cmd/goarticles/goarticles.go

.PHONY: run-api
run-api: build
	.build/goarticles-api --configfile=config.yml
