package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

type config struct {
	Server server `yaml:"server"`
	Mysql  mysql  `yaml:"mysql"`
}

type server struct {
	Host string `yaml:"host"`
}

type mysql struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Db       string `yaml:"db"`
}

var Config *config

func init() {
	yamlFile, err := os.ReadFile("./config.yaml")
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(yamlFile, &Config)
	if err != nil {
		panic(err)
	}
}
