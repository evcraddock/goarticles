package main

import (
	"github.com/evcraddock/goarticles/cli"
)

func main() {

	//TODO: get filename from flag
	var filename string

	cli.CreateOrUpdateArticle(filename)
}
