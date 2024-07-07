package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	calendarURL := os.Getenv("CALENDAR_URL")
	if calendarURL == "" {
		log.Fatal("CALENDAR_URL is not set in the .env file")
	}

	isDebug := os.Getenv("ENVIRONMENT") == "debug"

	c := cron.New(cron.WithLocation(time.FixedZone("UTC+5", 5*60*60)))

	_, err = c.AddFunc("20 10 * * 1-5", func() {
		fetchAndParseMeetingURL(calendarURL, isDebug)
	})
	if err != nil {
		log.Fatalf("Error scheduling cron job: %v", err)
	}

	c.Start()

	// keep the program running
	select {}
}

func fetchAndParseMeetingURL(calendarURL string, isDebug bool) {
	var reader io.Reader
	if isDebug {
		file, err := os.Open("index.html")
		if err != nil {
			log.Printf("Error opening index.html: %v", err)
			return
		}
		defer file.Close()

		reader = file
	} else {
		resp, err := http.Get(calendarURL)
		if err != nil {
			log.Printf("Error fetching URL: %v", err)
			return
		}
		defer resp.Body.Close()

		reader = resp.Body
	}

	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		log.Printf("Error parsing HTML: %v", err)
		return
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

	if meetingURL != "" {
		fmt.Printf("Video call URL: %s\n", meetingURL)
	} else {
		fmt.Println("No video call URL found.")
	}
}
