package models

import (
	log "github.com/sirupsen/logrus"
	"os"
	"time"
)

type Config struct {
	ServerAddress  string
	ServerPort     string
	LogLevel       log.Level
	DatabaseServer string
	DatabasePort   string
	DatabaseName   string
	TimeoutWait    time.Duration
}

func GetConfig() *Config {
	config := &Config{
		ServerAddress:  getConfigValue("ArticleServerAddress", "0.0.0.0").(string),
		ServerPort:     getConfigValue("ArticleServerPort", "8080").(string),
		LogLevel:       getConfigValue("ArticleServerLogLevel", log.DebugLevel).(log.Level),
		DatabaseServer: getConfigValue("ArticleServerDatabaseServer", "127.0.0.1").(string),
		DatabasePort:   getConfigValue("ArticleServerDatabasePort", "27017").(string),
		DatabaseName:   getConfigValue("ArticleServerDatabaseName", "articleDB").(string),
		TimeoutWait:    time.Second * time.Duration(getConfigValue("ArticleServerTimeoutWait", 15).(int)),
	}

	return config
}

func getConfigValue(variable string, defaultValue interface{}) interface{} {
	envVariable := os.Getenv(variable)

	if envVariable != "" {
		return envVariable
	}

	return defaultValue
}
