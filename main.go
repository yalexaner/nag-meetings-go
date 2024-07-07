package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/PuerkitoBio/goquery"
	"github.com/joho/godotenv"
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

	var reader io.Reader
	if isDebug {
		file, err := os.Open("index.html")
		if err != nil {
			log.Fatalf("Error opening index.html: %v", err)
		}
		defer file.Close()

		reader = file
	} else {
		resp, err := http.Get(calendarURL)
		if err != nil {
			log.Fatalf("Error fetching URL: %v", err)
		}
		defer resp.Body.Close()

		reader = resp.Body
	}

	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		log.Fatalf("Error parsing HTML: %v", err)
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
