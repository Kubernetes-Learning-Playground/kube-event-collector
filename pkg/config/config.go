package config

import (
	"fmt"
	"github.com/go-yaml/yaml"
	"io/ioutil"
	"log"
)

// Config 配置文件
type Config struct {
	KubeConfig  string `json:"kubeConfig" yaml:"kubeConfig"`
	FilterLevel string `json:"filterEventLevel" yaml:"filterEventLevel"`
	ElasticSearchEndpoint string `json:"elasticSearchEndpoint" yaml:"elasticSearchEndpoint"`
	Mode
	Sender `json:"sender" yaml:"sender"`
}

type Mode struct {
	Log        bool `json:"log" yaml:"log"`
	Prometheus bool `json:"prometheus", yaml:"prometheus"`
	Message    bool `json:"message" yaml:"message"`
	ElasticSearch bool `json:"elasticSearch" yaml:"elasticSearch"`
}

type Sender struct {
	Remote   string `yaml:"remote"`
	Port     int    `yaml:"port"`
	Email    string `yaml:"email"`
	Password string `yaml:"password"`
	Targets  string `yaml:"targets"`
}

func NewConfig() *Config {
	return &Config{}
}

func loadConfigFile(path string) []byte {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		log.Println(err)
		return nil
	}
	return b
}

func LoadConfig(path string) (*Config, error) {
	config := NewConfig()
	if b := loadConfigFile(path); b != nil {

		err := yaml.Unmarshal(b, config)
		if err != nil {
			return nil, err
		}
		return config, err
	} else {
		return nil, fmt.Errorf("load config1 file error...")
	}

}
