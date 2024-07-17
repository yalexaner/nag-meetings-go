package main

import (
	"database/sql"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/PuerkitoBio/goquery"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
	"github.com/robfig/cron/v3"
	"github.com/yalexaner/nag-meetings-go/messages"
)

var (
	db  *sql.DB
	bot *tgbotapi.BotAPI
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	calendarURL := os.Getenv("CALENDAR_URL")
	if calendarURL == "" {
		log.Fatal("CALENDAR_URL is not set in the .env file")
	}

	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	if botToken == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN is not set in the .env file")
	}

	workingDirectory := os.Getenv("WORKING_DIRECTORY")
	if workingDirectory == "" {
		log.Fatal("WORKING_DIRECTORY is not set in the .env file")
	}

	isDebug := os.Getenv("ENVIRONMENT") == "debug"

	db, err = sql.Open("sqlite3", workingDirectory+"subscribers.db")
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS subscribers (user_id INTEGER PRIMARY KEY)`)
	if err != nil {
		log.Fatalf("Error creating table: %v", err)
	}

	bot, err = tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Fatalf("Error initializing bot: %v", err)
	}

	bot.Debug = isDebug

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Fatalf("Error getting updates channel: %v", err)
	}

	// Start cron job
	c := cron.New(cron.WithLocation(time.FixedZone("UTC+5", 5*60*60)))

	_, err = c.AddFunc("20 10 * * 1-5", func() {
		meetingURL := fetchAndParseMeetingURL(calendarURL, isDebug)
		if meetingURL != "" {
			sendMeetingURLToSubscribers(meetingURL)
		}
	})
	if err != nil {
		log.Fatalf("Error scheduling cron job: %v", err)
	}

	c.Start()

	handleUpdates(updates)
}

func handleUpdates(updates tgbotapi.UpdatesChannel) {
	for update := range updates {
		if update.Message == nil {
			continue
		}

		switch update.Message.Command() {
		case "subscribe":
			handleSubscribe(update.Message.Chat.ID)
		case "unsubscribe":
			handleUnsubscribe(update.Message.Chat.ID)
		default:
			sendMessage(update.Message.Chat.ID, messages.UnknownCommand)
		}
	}
}

func handleSubscribe(chatID int64) {
	_, err := db.Exec("INSERT OR IGNORE INTO subscribers (user_id) VALUES (?)", chatID)
	if err != nil {
		log.Printf("Error subscribing user: %v", err)
		sendMessage(chatID, messages.ErrorSubscribing)
		return
	}
	sendMessage(chatID, messages.Subscribed)
}

func handleUnsubscribe(chatID int64) {
	_, err := db.Exec("DELETE FROM subscribers WHERE user_id = ?", chatID)
	if err != nil {
		log.Printf("Error unsubscribing user: %v", err)
		sendMessage(chatID, messages.ErrorUnsubscribing)
		return
	}
	sendMessage(chatID, messages.Unsubscribed)
}

func sendMessage(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	_, err := bot.Send(msg)
	if err != nil {
		log.Printf("Error sending message: %v", err)
	}
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

func sendMeetingURLToSubscribers(meetingURL string) {
	rows, err := db.Query("SELECT user_id FROM subscribers")
	if err != nil {
		log.Printf("Error querying subscribers: %v", err)
		return
	}
	defer rows.Close()

	var subscribers []int64
	for rows.Next() {
		var userID int64
		err := rows.Scan(&userID)
		if err != nil {
			log.Printf("Error scanning user ID: %v", err)
			continue
		}
		subscribers = append(subscribers, userID)
	}

	for _, userID := range subscribers {
		sendMessage(userID, meetingURL)
		time.Sleep(time.Duration(1000/30) * time.Millisecond)
	}
}
