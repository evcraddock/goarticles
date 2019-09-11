package main

import (
	"fmt"
	"os"

	"github.com/evcraddock/goarticles/pkg/cmd"
)

func main() {
	command := cmd.NewDefaultCommand()
	if err := command.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	//	log.SetFormatter(&log.JSONFormatter{})
	//	log.SetOutput(os.Stdout)
	//
	//	configFile := flag.String("configfile", "", "yaml configuration file (optional)")
	//	filesToProcess := flag.String("files", "", "files or folders to process")
	//	flag.Parse()
	//
	//	var config *configs.ClientConfiguration
	//	var err error
	//
	//	if config, err = configs.LoadCliConfig(*configFile); err != nil {
	//		log.Error(err.Error())
	//		return
	//	}
	//
	//	log.SetLevel(log.InfoLevel)
	//
	//	articleService := cli.NewArticleImporter(*config)
	//	articleService.CreateOrUpdateArticle(*filesToProcess)
}
