package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type ParserConfig struct {
	EventSelector   string // CSS selector for event containers
	TitleSelector   string // CSS selector for event title within event
	MeetingName     string // Meeting name to match
	LinkSelector    string // CSS selector for meeting links
	LinkAttribute   string // HTML attribute to extract (usually "href")
}

type Config struct {
	CalendarURL      string
	TelegramBotToken string
	WorkingDirectory string
	IsDebug          bool
	Parser           ParserConfig
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
		Parser: ParserConfig{
			EventSelector: getEnvOrDefault("PARSER_EVENT_SELECTOR", ".b-content-event"),
			TitleSelector: getEnvOrDefault("PARSER_TITLE_SELECTOR", "h1"),
			MeetingName:   getEnvOrDefault("PARSER_MEETING_NAME", "STB Daily Meeting"),
			LinkSelector:  getEnvOrDefault("PARSER_LINK_SELECTOR", ".e-description a"),
			LinkAttribute: getEnvOrDefault("PARSER_LINK_ATTRIBUTE", "href"),
		},
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

func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
