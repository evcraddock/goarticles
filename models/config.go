package models

import (
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

//Config configuration settings
type Config struct {
	ServerAddress  string
	ServerPort     string
	LogLevel       log.Level
	DatabaseServer string
	DatabasePort   string
	DatabaseName   string
	TimeoutWait    time.Duration
	AuthDomain     string
	Audience       string
}

//GetConfig get configuration value from env variable or return default
func GetConfig() *Config {
	config := &Config{
		ServerAddress:  getConfigValue("ArticleServerAddress", "0.0.0.0").(string),
		ServerPort:     getConfigValue("ArticleServerPort", "8080").(string),
		LogLevel:       getConfigValue("ArticleServerLogLevel", log.DebugLevel).(log.Level),
		DatabaseServer: getConfigValue("ArticleServerDatabaseServer", "127.0.0.1").(string),
		DatabasePort:   getConfigValue("ArticleServerDatabasePort", "27017").(string),
		DatabaseName:   getConfigValue("ArticleServerDatabaseName", "articleDB").(string),
		TimeoutWait:    time.Second * time.Duration(getConfigValue("ArticleServerTimeoutWait", 15).(int)),
		AuthDomain:     getConfigValue("ArticleServiceAuthDomain", "erikvan.auth0.com").(string),
		Audience:       getConfigValue("ArticleServiceAudience", "https://api.erikvancraddock.com").(string),
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
