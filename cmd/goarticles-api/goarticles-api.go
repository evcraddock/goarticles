package main

import (
	log "github.com/sirupsen/logrus"

	"os"

	"flag"

	"github.com/evcraddock/goarticles/api"
	"github.com/evcraddock/goarticles/configs"
)

func main() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)

	configFile := flag.String("configfile", "", "yaml configuration file (optional)")
	flag.Parse()

	log.Debug("Loading configuration")
	config, err := configs.LoadEnvironmentVariables()
	if err != nil {
		log.Debug("no environment variables found, checking config file")
		config, err = configs.LoadConfigFile(*configFile)
		if err != nil {
			log.Error(err.Error())
			panic(err)
		}
	}

	loglevel := setLogLevel(config.Server.LogLevel)
	log.Debugf("LogLevel: %v", loglevel)
	log.SetLevel(loglevel)

	api.NewServer(config)

}

func setLogLevel(logLevel string) log.Level {
	switch logLevel {
	case "debug":
		return log.DebugLevel
	case "error":
		return log.ErrorLevel
	default:
		return log.InfoLevel
	}
}
