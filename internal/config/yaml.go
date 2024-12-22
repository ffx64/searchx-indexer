package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type config struct {
	Database struct {
		Host           string `yaml:"host"`
		Port           int32  `yaml:"port"`
		Username       string `yaml:"username"`
		Password       string `yaml:"password"`
		Dbname         string `yaml:"dbname"`
		Sslmode        string `yaml:"sslmode"`
		MaxConnections int    `yaml:"max_connections"`
		Timeout        int    `yaml:"timeout"`
	} `yaml:"database"`
}

func LoadConfig(filename string) (*config, error) {
	var config config
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
