package main

import (
	"flag"
	"os"

	"github.com/evcraddock/goarticles/proxy"

	log "github.com/sirupsen/logrus"

	"github.com/evcraddock/goarticles/configs"
)

func main() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)

	configFile := flag.String("configfile", "", "yaml configuration file (optional)")
	flag.Parse()

	log.Info("Loading configuration from environment variables")
	config, err := configs.LoadProxyEnvironmentVariables()
	if err != nil {
		config, err = configs.LoadProxyConfigFile(*configFile)
		if err != nil {
			log.Error(err.Error())
			panic(err)
		}
	}

	loglevel := setLogLevel(config.Server.LogLevel)
	log.Infof("LogLevel: %v", loglevel)
	log.SetLevel(loglevel)

	reverseProxy := proxy.NewServer(config)
	reverseProxy.Start()
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
