package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type config struct {
	Server        server              `yaml:"server"`
	Db            db                  `yaml:"db"`
	Elasticsearch elasticsearchConfig `yaml:"elasticsearch"`
}

type server struct {
	Port string `yaml:"port"`
}

type db struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Db       string `yaml:"db"`
}

// elasticsearchConfig ES配置
type elasticsearchConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

var Config *config

func LoadConfig(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("无法读取配置文件: %w", err)
	}

	if err := yaml.Unmarshal(data, &Config); err != nil {
		return fmt.Errorf("无法解析配置文件: %w", err)
	}

	return nil
}
