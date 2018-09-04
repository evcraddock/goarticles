package configs

import (
	"io/ioutil"

	log "github.com/sirupsen/logrus"

	"os"

	"fmt"

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
	authConfig := &AuthConfig{}
	errors := make([]error, 0)

	if url, exists := os.LookupEnv("CLI_AUTH_URL"); exists {
		authConfig.URL = url
	} else {
		errors = append(errors, fmt.Errorf("env variable CLI_AUTH_URL does not exist"))
	}

	if grantType, exists := os.LookupEnv("CLI_GRANT_TYPE"); exists {
		authConfig.GrantType = grantType
	} else {
		errors = append(errors, fmt.Errorf("env variable CLI_GRANT_TYPE does not exist"))
	}

	if clientId, exists := os.LookupEnv("CLI_CLIENT_ID"); exists {
		authConfig.ClientID = clientId
	} else {
		errors = append(errors, fmt.Errorf("env variable CLI_CLIENT_IDdoes not exist"))
	}

	if clientSecret, exists := os.LookupEnv("CLI_CLIENT_SECRET"); exists {
		authConfig.ClientSecret = clientSecret
	} else {
		errors = append(errors, fmt.Errorf("env variable CLI_CLIENT_SECRET does not exist"))
	}

	if audience, exists := os.LookupEnv("CLI_AUTH_AUDIENCE"); exists {
		authConfig.Audience = audience
	} else {
		errors = append(errors, fmt.Errorf("env variable CLI_AUTH_AUDIENCE does not exist"))
	}

	cliConfig := &ClientConfiguration{
		Auth: *authConfig,
	}

	if apiUrl, exists := os.LookupEnv("CLI_API_URL"); exists {
		cliConfig.URL = apiUrl
	} else {
		errors = append(errors, fmt.Errorf("env variable CLI_API_URL does not exist"))
	}

	if len(errors) > 0 {
		for _, v := range errors {
			log.Error(v.(error).Error())
		}

		return nil, fmt.Errorf("could not load env variables")
	}

	return cliConfig, nil
}
