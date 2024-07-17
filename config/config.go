package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	CalendarURL      string
	TelegramBotToken string
	WorkingDirectory string
	IsDebug          bool
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	config := &Config{
		CalendarURL:      getEnvOrFatal("CALENDAR_URL"),
		TelegramBotToken: getEnvOrFatal("TELEGRAM_BOT_TOKEN"),
		WorkingDirectory: getEnvOrFatal("WORKING_DIRECTORY"),
		IsDebug:          os.Getenv("ENVIRONMENT") == "debug",
	}

	return config
}

func getEnvOrFatal(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("%s is not set in the .env file", key)
	}
	return value
}
