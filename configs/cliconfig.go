package configs

import (
	log "github.com/sirupsen/logrus"
	//"gopkg.in/yaml.v2"
	"io/ioutil"

	"os"

	"gopkg.in/yaml.v2"
)

//ClientConfiguration top level client config object
type ClientConfiguration struct {
	URL  string     `yaml:"api-url"`
	Auth AuthConfig `yaml:"auth"`
}

//AuthConfig authorization config data
type AuthConfig struct {
	URL          string `yaml:"auth-url"`
	GrantType    string `yaml:"grant_type"`
	ClientID     string `yaml:"client_id"`
	ClientSecret string `yaml:"client_secret"`
	Audience     string `yaml:"audience"`
}

//LoadCliConfig load configuration for client
func LoadCliConfig(configFile string) (*ClientConfiguration, error) {
	if configFile != "" {
		log.Debug("Loading configuration from config file")
		return LoadCliConfigFromFile(configFile)
	}

	log.Debug("Loading configuration from environment variables")
	return LoadCliConfigFromEnvVariables()
}

//LoadCliConfigFromFile loading configuration from yaml file
func LoadCliConfigFromFile(filename string) (*ClientConfiguration, error) {
	b, err := ioutil.ReadFile(filename)

	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	var config ClientConfiguration
	err = yaml.Unmarshal(b, &config)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	return &config, nil
}

//LoadCliConfigFromEnvVariables load configuration from env variables
func LoadCliConfigFromEnvVariables() (*ClientConfiguration, error) {
	return &ClientConfiguration{
		URL: os.Getenv("CLI_API_URL"),
		Auth: AuthConfig{
			URL:          os.Getenv("CLI_AUTH_URL"),
			GrantType:    os.Getenv("CLI_GRANT_TYPE"),
			ClientID:     os.Getenv("CLI_CLIENT_ID"),
			ClientSecret: os.Getenv("CLI_CLIENT_SECRET"),
			Audience:     os.Getenv("CLI_AUTH_AUDIENCE"),
		},
	}, nil
}
