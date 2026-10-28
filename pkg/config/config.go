package config

import (
	"os"

	"go.yaml.in/yaml/v3"
)

type Config struct {
	Name string `yaml:"Name"`
	Port string `yaml:"Port"`
	Host string `yaml:"Host"`
}

var Conf Config = Config{}

func init() {
	data, err := os.ReadFile("./config.yaml")
	if err != nil {
		panic(err)
		return
	}
	yaml.Unmarshal(data, &Conf)

}
