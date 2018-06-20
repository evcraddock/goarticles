package models

import (
	"os"
	"time"

	"io/ioutil"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

//Configuration top level config object
type Configuration struct {
	Server         ServerConfiguration         `yaml:"server"`
	Database       DatabaseConfiguration       `yaml:"database"`
	Authentication AuthenticationConfiguration `yaml:"authentication"`
}

//ServerConfiguration server config data
type ServerConfiguration struct {
	Address  string `yaml:"address"`
	Port     string `yaml:"port"`
	LogLevel string `yaml:"loglevel"`
}

//DatabaseConfiguration database config data
type DatabaseConfiguration struct {
	Address      string        `yaml:"address"`
	Port         string        `yaml:"port"`
	DatabaseName string        `yaml:"databasename"`
	Timeout      time.Duration `yaml:"timeout"`
}

//AuthenticationConfiguration authentication config data
type AuthenticationConfiguration struct {
	Domain   string `yaml:"domain"`
	Audience string `yaml:"audience"`
}

//LoadConfig from file
func LoadConfig(filename string) (*Configuration, error) {

	b, err := ioutil.ReadFile(filename)

	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	var config Configuration

	err = yaml.Unmarshal(b, &config)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	return &config, nil
}

//LoadEnvironmentVariables load from env variables
//TODO: return error if env variables are empty
func LoadEnvironmentVariables() (*Configuration, error) {
	timeout, err := time.ParseDuration(os.Getenv("GOA_DB_TIMEOUT"))
	if err != nil {
		return nil, err
	}

	return &Configuration{
		ServerConfiguration{
			Address:  os.Getenv("GOA_SERVER_ADDRESS"),
			Port:     os.Getenv("GOA_SERVER_PORT"),
			LogLevel: os.Getenv("GOA_LOG_LEVEL"),
		},
		DatabaseConfiguration{
			Address:      os.Getenv("GOA_DB_ADDRESS"),
			Port:         os.Getenv("GOA_DB_PORT"),
			DatabaseName: os.Getenv("GOA_DB_DATABASENAME"),
			Timeout:      timeout,
		},
		AuthenticationConfiguration{
			Domain:   os.Getenv("GOA_AUTH_DOMAIN"),
			Audience: os.Getenv("GOA_AUTH_AUDIENCE"),
		},
	}, nil
}
