package main

import (
	"flag"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/evcraddock/goarticles/internal/cli"
	"github.com/evcraddock/goarticles/internal/configs"
)

func main() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)

	configFile := flag.String("configfile", "", "yaml configuration file (optional)")
	filesToProcess := flag.String("files", "", "files or folders to process")
	flag.Parse()

	var config *configs.ClientConfiguration
	var err error

	if config, err = configs.LoadCliConfig(*configFile); err != nil {
		log.Error(err.Error())
		return
	}

	log.SetLevel(log.InfoLevel)

	articleService := cli.NewArticleImporter(*config)
	articleService.CreateOrUpdateArticle(*filesToProcess)
}
