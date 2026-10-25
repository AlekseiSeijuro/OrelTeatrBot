package config

import "os"

type Config struct {
	DbConnectionString string
	BotToken           string
}

func ReadConfigFromEnv() (*Config, error) {
	return &Config{
		os.Getenv("DATABASE_URL"),
		os.Getenv("BOT_TOKEN")}, nil
}
