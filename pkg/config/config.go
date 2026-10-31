package config

import (
	"os"
	"sync"

	"go.yaml.in/yaml/v3"
)

type Config struct {
	Name string `yaml:"Name"`
	Port string `yaml:"Port"`
	Host string `yaml:"Host"`

	DataFolder   string `yaml:"DataFolder"`
	OneFileName  string `yaml:"OneFileName"`
	IsOverride   bool   `yaml:"IsOverride"`
	IsMultiFiles bool   `yaml:"IsMultiFiles"`
}

var Conf *Config

func init() {
	sync.OnceFunc(func() {
		data, err := os.ReadFile("./config.yaml")
		if err != nil {
			panic(err)
		}
		Conf = &Config{
			Name:         "Metrics-Server",
			Port:         "50051",
			IsMultiFiles: false,
			IsOverride:   false,
			DataFolder:   "data/",
			OneFileName:  "MainMetrics.txt",
		}

		yaml.Unmarshal(data, Conf)
	})()
}
