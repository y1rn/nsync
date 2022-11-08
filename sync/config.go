package sync

import (
	"os"

	"gopkg.in/yaml.v3"
)

const (
	KEY_NAME    = "name"
	KEY_SERVER  = "server"
	KEY_PORT    = "port"
	KEY_TYPE    = "type"
	KEY_GROUP   = "proxy-groups"
	KEY_PROXIES = "proxies"
)

type ProxyGroup struct {
	Name    string   `yaml:"name"`
	Proxies []string `yaml:"proxies"`
	Type    string   `yaml:"type"`
	Url     string   `yaml:"url"`
}
type Config struct {
	TemplatePath string           `yaml:"temp"`
	OutPut       string           `yaml:"output"`
	Group        ProxyGroup       `yaml:"proxy-group"`
	Proxies      []map[string]any `yaml:"proxies"`
}

func (c *Config) LoadTemplate() (map[string]any, error) {
	data, err := os.ReadFile(c.TemplatePath)
	if err != nil {
		return nil, err
	}
	temp := make(map[string]any)
	err = yaml.Unmarshal(data, &temp)
	if err != nil {
		return nil, err
	}
	return temp, nil
}

func loadConfig(configPath string) ([]Config, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}
	var configs []Config
	err = yaml.Unmarshal(data, &configs)
	if err != nil {
		return nil, err
	}
	return configs, nil
}
