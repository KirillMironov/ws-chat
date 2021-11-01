package config

import (
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Port string `default:"8080" envconfig:"PORT"`
}

func InitConfig() (*Config, error) {
	var config Config
	return &config, envconfig.Process("", &config)
}
