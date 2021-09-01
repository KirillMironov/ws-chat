package config

import (
	"log"
	"os"
)

type Config struct {
	Port string
}

func InitConfig() *Config {
	var config Config

	config.Port = getEnv("PORT")
	return &config
}

func getEnv(key string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		log.Fatalf("Environment variable %s was not found", key)
		return ""
	}
	return value
}
