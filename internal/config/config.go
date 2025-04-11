package config

import "os"

type Config struct {
	Port string `mapstructure:"PORT"`
}

func LoadConfig() *Config {
	return &Config{
		Port: getEnv("PORT", "8080"),
	}
}

func getEnv(key string, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}
