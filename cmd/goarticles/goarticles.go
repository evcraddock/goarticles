package main

import (
	"flag"
	"os"

	"github.com/evcraddock/goarticles/cli"
	"github.com/evcraddock/goarticles/configs"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)

	configFile := flag.String("configfile", "", "yaml configuration file (optional)")
	filesToProcess := flag.String("files", "", "files or folders to process")
	flag.Parse()

	log.Info("Loading configuration from environment variables")
	config, err := configs.LoadCliConfigFromEnvVariables()
	if err != nil {
		config, err = configs.LoadCliConfig(*configFile)
		if err != nil {
			log.Error(err.Error())
			panic(err)
		}
	}

	log.SetLevel(log.InfoLevel)

	articleService := cli.NewImportArticleService(*config)
	articleService.CreateOrUpdateArticle(*filesToProcess)
}
