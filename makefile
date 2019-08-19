.PHONY: build
build:
	GOOS=linux GOARCH=amd64 go build -o .build/goarticles-api cmd/goarticles-api/goarticles-api.go

.PHONY: run
run:
	.build/goarticles-api --configfile=config.yml
