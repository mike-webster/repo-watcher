package env

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	yaml "gopkg.in/yaml.v1"
)

var curConfig *Config

// Watcher is a configuration for a source (repo) and destination (slack)
type Watcher struct {
	Repo    string `yaml:"repo"`
	Webhook string `yaml:"webhook"`
}

type Watchers []Watcher

func (w Watchers) Select(repo string) *Watcher {
	for _, r := range w {
		if strings.ToLower(r.Repo) == strings.ToLower(repo) {
			return &r
		}
	}
	return nil
}

func (w Watchers) ToString() string {
	ret := []string{}
	for _, r := range w {
		ret = append(ret, r.Repo)
	}
	return strings.Join(ret, ",")
}

// Config contains all the application's configured values
type Config struct {
	Port            int      `yaml:"port"`
	APIToken        string   `yaml:"token"`
	BaseURLTemplate string   `yaml:"base_url_template"`
	OrgName         string   `yaml:"org_name"`
	EventEndpoint   string   `yaml:"event_endpoint"`
	UserEndpoint    string   `yaml:"user_endpoint"`
	RefreshTimer    int      `yaml:"refresh_seconds"`
	RepoHost        string   `yaml:"repo_host"`
	RepoToWatch     string   `yaml:"repo_to_watch"`
	LogLevel        string   `yaml:"log_level"`
	SlackWebhook    string   `yaml:"slack_webhook"`
	RunType         string   `yaml:"run_type"`
	Watchers        Watchers `yaml:"watchers"`
}

func (c *Config) BaseURL() string {
	return fmt.Sprintf(c.BaseURLTemplate, c.RepoHost)
}

var (
	// AppConfigFile is a relative file path
	AppConfigFile = "app.yaml"
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

	var configs map[string]Config
	if err = yaml.Unmarshal(f, &configs); err != nil {
		panic(err)
	}

	env := os.Getenv("GO_ENV")
	fmt.Println("env: ", env)
	dev := configs[env]
	envToken := os.Getenv("API_TOKEN")
	if len(envToken) > 0 {
		dev.APIToken = envToken
	}

	return &dev
}
