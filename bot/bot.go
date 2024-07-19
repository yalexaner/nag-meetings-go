package bot

import (
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/yalexaner/nag-meetings-go/database"
	"github.com/yalexaner/nag-meetings-go/messages"
)

type Bot struct {
	api *tgbotapi.BotAPI
	db  *database.Database
}

func NewBot(token string, db *database.Database, debug bool) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	api.Debug = debug

	log.Printf("Authorized on account %s", api.Self.UserName)

	return &Bot{api: api, db: db}, nil
}

func (b *Bot) Start() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := b.api.GetUpdatesChan(u)
	if err != nil {
		log.Fatalf("Error getting updates channel: %v", err)
	}

	b.handleUpdates(updates)
}

func (b *Bot) handleUpdates(updates tgbotapi.UpdatesChannel) {
	for update := range updates {
		if update.Message == nil {
			continue
		}

		switch update.Message.Command() {
		case "start":
			b.sendMessage(update.Message.Chat.ID, messages.Start)
		case "subscribe":
			b.handleSubscribe(update.Message.Chat.ID)
		case "unsubscribe":
			b.handleUnsubscribe(update.Message.Chat.ID)
		default:
			b.sendMessage(update.Message.Chat.ID, messages.UnknownCommand)
		}
	}
}

func (b *Bot) handleStart(chatId int64) {
	if err := b.db.AddNewUser(chatId); err != nil {
		log.Printf("Error adding new user: %v", err)
		return
	}

	b.sendMessage(chatId, messages.Start)
}

func (b *Bot) handleSubscribe(chatID int64) {
	if err := b.db.Subscribe(chatID); err != nil {
		log.Printf("Error subscribing user: %v", err)
		b.sendMessage(chatID, messages.ErrorSubscribing)
		return
	}

	b.sendMessage(chatID, messages.Subscribed)
}

func (b *Bot) handleUnsubscribe(chatID int64) {
	if err := b.db.Unsubscribe(chatID); err != nil {
		log.Printf("Error unsubscribing user: %v", err)
		b.sendMessage(chatID, messages.ErrorUnsubscribing)
		return
	}

	b.sendMessage(chatID, messages.Unsubscribed)
}

func (b *Bot) sendMessage(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	_, err := b.api.Send(msg)
	if err != nil {
		log.Printf("Error sending message: %v", err)
	}
}

func (b *Bot) SendMeetingURLToSubscribers(meetingURL string) {
	subscribers, err := b.db.GetSubscribers()
	if err != nil {
		log.Printf("Error querying subscribers: %v", err)
		return
	}

	for _, userID := range subscribers {
		b.sendMessage(userID, meetingURL)
		time.Sleep(time.Duration(1000/30) * time.Millisecond)
	}
}
