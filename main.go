package main

import (
	"log"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/yalexaner/nag-meetings-go/bot"
	"github.com/yalexaner/nag-meetings-go/config"
	"github.com/yalexaner/nag-meetings-go/database"
	"github.com/yalexaner/nag-meetings-go/parser"
)

func main() {
	cfg := config.LoadConfig()

	db, err := database.NewDatabase(cfg.WorkingDirectory + "subscribers.db")
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()

	tgbot, err := bot.NewBot(cfg.TelegramBotToken, db, cfg.IsDebug)
	if err != nil {
		log.Fatalf("Error initializing bot: %v", err)
	}

	c := cron.New(cron.WithLocation(time.FixedZone("UTC+5", 5*60*60)))

	_, err = c.AddFunc("20 10 * * 1-5", func() {
		meetingURL := parser.FetchMeetingURL(cfg.CalendarURL, cfg.IsDebug)
		if meetingURL != "" {
			tgbot.SendMeetingURLToSubscribers(meetingURL)
		}
	})
	if err != nil {
		log.Fatalf("Error scheduling cron job: %v", err)
	}

	c.Start()

	tgbot.Start()
}
