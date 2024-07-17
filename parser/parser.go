package parser

import (
	"io"
	"log"
	"net/http"
	"os"

	"github.com/PuerkitoBio/goquery"
)

func FetchMeetingURL(calendarURL string, isDebug bool) string {
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
