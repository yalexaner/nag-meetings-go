package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/robfig/cron/v3"
	"github.com/yalexaner/nag-meetings-go/bot"
	"github.com/yalexaner/nag-meetings-go/config"
	"github.com/yalexaner/nag-meetings-go/database"
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
		meetingURL := fetchAndParseMeetingURL(cfg.CalendarURL, cfg.IsDebug)
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

func fetchAndParseMeetingURL(calendarURL string, isDebug bool) string {
	var reader io.Reader
	if isDebug {
		file, err := os.Open("index.html")
		if err != nil {
			log.Printf("Error opening index.html: %v", err)
			return ""
		}
		defer file.Close()

		reader = file
	} else {
		resp, err := http.Get(calendarURL)
		if err != nil {
			log.Printf("Error fetching URL: %v", err)
			return ""
		}
		defer resp.Body.Close()

		reader = resp.Body
	}

	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		log.Printf("Error parsing HTML: %v", err)
		return ""
	}

	var meetingURL string
	doc.Find(".b-content-event").EachWithBreak(func(i int, s *goquery.Selection) bool {
		if s.Find("h1").Text() != "STB Daily Meeting" {
			return true
		}

		s.Find(".e-description a").Each(func(i int, a *goquery.Selection) {
			if meetingURL == "" {
				meetingURL, _ = a.Attr("href")
			}
		})

		return false
	})

	return meetingURL
}
