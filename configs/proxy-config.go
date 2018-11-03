package configs

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

//ProxyServerConfiguration server config data
type ProxyServerConfiguration struct {
	Port        string        `yaml:"port"`
	LogLevel    string        `yaml:"loglevel"`
	Timeout     time.Duration `yaml:"timeout"`
	StaticFiles string        `yaml:"staticfiles"`
	ForwardAPI  string        `yaml:"forwardapi"`
	WhiteList   string        `yaml:"whitelist"`
}

//ProxyConfiguration top level config object
type ProxyConfiguration struct {
	Server ProxyServerConfiguration `yaml:"server"`
	Auth   AuthConfig               `yaml:"auth"`
}

//LoadProxyConfigFile load from file
func LoadProxyConfigFile(filename string) (*ProxyConfiguration, error) {

	b, err := ioutil.ReadFile(filename)

	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	var config ProxyConfiguration

	err = yaml.Unmarshal(b, &config)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	return &config, nil
}

//LoadProxyEnvironmentVariables load from env variables
func LoadProxyEnvironmentVariables() (*ProxyConfiguration, error) {
	timeout, err := time.ParseDuration(os.Getenv("GOAP_TIMEOUT"))
	if err != nil {
		return nil, err
	}

	authConfig := &AuthConfig{}
	errors := make([]error, 0)

	if url, exists := os.LookupEnv("GOAP_AUTH_URL"); exists {
		authConfig.URL = url
	} else {
		errors = append(errors, fmt.Errorf("env variable GOAP_AUTH_URL does not exist"))
	}

	if grantType, exists := os.LookupEnv("GOAP_GRANT_TYPE"); exists {
		authConfig.GrantType = grantType
	} else {
		errors = append(errors, fmt.Errorf("env variable GOAP_GRANT_TYPE does not exist"))
	}

	if clientID, exists := os.LookupEnv("GOAP_CLIENT_ID"); exists {
		authConfig.ClientID = clientID
	} else {
		errors = append(errors, fmt.Errorf("env variable GOAP_CLIENT_IDdoes not exist"))
	}

	if clientSecret, exists := os.LookupEnv("GOAP_CLIENT_SECRET"); exists {
		authConfig.ClientSecret = clientSecret
	} else {
		errors = append(errors, fmt.Errorf("env variable GOAP_CLIENT_SECRET does not exist"))
	}

	if audience, exists := os.LookupEnv("GOAP_AUTH_AUDIENCE"); exists {
		authConfig.Audience = audience
	} else {
		errors = append(errors, fmt.Errorf("env variable GOAP_AUTH_AUDIENCE does not exist"))
	}

	return &ProxyConfiguration{
		ProxyServerConfiguration{
			Port:        os.Getenv("GOAP_SERVER_PORT"),
			LogLevel:    os.Getenv("GOAP_LOG_LEVEL"),
			Timeout:     timeout,
			StaticFiles: os.Getenv("GOAP_STATIC_FILES"),
			ForwardAPI:  os.Getenv("GOAP_FORWARD_API"),
		},
		*authConfig,
	}, nil
}
