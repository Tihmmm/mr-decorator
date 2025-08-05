package config

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/goccy/go-yaml"
	"os"
)

type ServerConfig struct {
	Port             string   `yaml:"port" default:"3000"`
	RateLimit        int      `yaml:"rate_limit" default:"3"`
	ApiKey           string   `yaml:"api_key"`
	SupportedFormats []string `yaml:"supported_formats"`
}

func NewConfig(path string) (ServerConfig, error) {
	configBytes, err := os.ReadFile(path)
	if err != nil {
		return ServerConfig{}, errors.New(fmt.Sprintf("Error reading config.yml: %s\n", err))
	}

	var cfg ServerConfig
	buf := bytes.NewBuffer(configBytes)
	dec := yaml.NewDecoder(buf)
	if err := dec.Decode(&cfg); err != nil {
		return ServerConfig{}, errors.New(fmt.Sprintf("Error parsing config.yml: %s\n", err))
	}

	return cfg, nil
}
