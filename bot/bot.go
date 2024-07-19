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
		if update.Message != nil && update.Message.IsCommand() {

			if update.Message.Command() == "start" {
				b.handleStartCommand(update.Message.Chat.ID)
				continue
			}

			isAuthorized := b.checkIsAuthorized(update.Message.Chat.ID)
			if !isAuthorized {
				b.sendMessage(update.Message.Chat.ID, messages.NotAuthorized)
				continue
			}

			switch update.Message.Command() {
			case "admin":
				b.handleAdminCommand(update.Message.Chat.ID)
			case "subscribe":
				b.handleSubscribeCommand(update.Message.Chat.ID)
			case "unsubscribe":
				b.handleUnsubscribeCommand(update.Message.Chat.ID)
			default:
				b.sendMessage(update.Message.Chat.ID, messages.UnknownCommand)
			}
		} else if update.CallbackQuery != nil {
			b.handleCallbackQuery(update.CallbackQuery)
		}
	}
}

func (b *Bot) checkIsAuthorized(chatId int64) bool {
	isAuthorized, err := b.db.IsAuthorized(chatId)
	if err != nil {
		log.Printf("Error checking if user is authorized: %v", err)
		return false
	}

	if isAuthorized == 1 {
		return true
	} else {
		return false
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
