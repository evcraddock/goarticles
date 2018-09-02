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

	testfilename := "/home/erik/code/src/github.com/evcraddock/erikvancraddock-articles/how-to-use-goals/how-to-use-goals.md"

	configFile := flag.String("configfile", "", "yaml configuration file (optional)")
	filesToProcess := flag.String("files", testfilename, "files or folders to process")
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
