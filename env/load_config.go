package env

import (
	"io/ioutil"
	"path/filepath"

	yaml "gopkg.in/yaml.v1"
)

var curConfig *Config

// Config contains all the application's configured values
type Config struct {
	Port          int    `yaml:"port"`
	APIToken      string `yaml:"token"`
	BaseURL       string `yaml:"base_url"`
	OrgName       string `yaml:"org_name"`
	EventEndpoint string `yaml:"event_endpoint"`
}

var (
	// AppConfigFile is a relative file path
	AppConfigFile = "app.yml"
	// BasePath is an absolute path to the directory holding the configuration data.
	BasePath string
)

// GetConfig retrieves the application's configuration values.
func GetConfig() *Config {
	if curConfig == nil {
		curConfig = loadAppConfig()
	}

	return curConfig
}

func loadAppConfig() *Config {
	path, err := filepath.Abs(AppConfigFile)
	if err != nil {
		panic(err)
	}

	f, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	BasePath, err = filepath.Abs(AppConfigFile)
	if err != nil {
		panic(err)
	}
	BasePath = filepath.Dir(BasePath)

	var config Config
	if err = yaml.Unmarshal(f, config); err != nil {
		panic(err)
	}

	return &config
}
