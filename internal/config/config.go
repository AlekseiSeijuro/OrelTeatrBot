package config

import "os"

type Config struct {
	DbConnectionString string
}

func ReadConfigFromEnv() (*Config, error) {
	return &Config{os.Getenv("DATABASE_URL")}, nil
}
