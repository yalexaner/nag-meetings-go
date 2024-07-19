package bot

import (
	"log"

	"github.com/yalexaner/nag-meetings-go/messages"
)

func (b *Bot) handleStartCommand(chatId int64) {
	if err := b.db.AddNewUser(chatId); err != nil {
		log.Printf("Error adding new user: %v", err)
		return
	}

	b.sendMessage(chatId, messages.Start)
}

func (b *Bot) handleSubscribeCommand(chatID int64) {
	if err := b.db.Subscribe(chatID); err != nil {
		log.Printf("Error subscribing user: %v", err)
		b.sendMessage(chatID, messages.ErrorSubscribing)
		return
	}

	b.sendMessage(chatID, messages.Subscribed)
}

func (b *Bot) handleUnsubscribeCommand(chatID int64) {
	if err := b.db.Unsubscribe(chatID); err != nil {
		log.Printf("Error unsubscribing user: %v", err)
		b.sendMessage(chatID, messages.ErrorUnsubscribing)
		return
	}

	b.sendMessage(chatID, messages.Unsubscribed)
}

func (b *Bot) handleAdminCommand(chatId int64) {
	isAdmin, err := b.db.IsAdmin(chatId)
	if err != nil {
		log.Printf("Error checking if user is admin: %v", err)
		b.sendMessage(chatId, messages.UnknownError)
		return
	}

	if isAdmin != 1 {
		b.sendMessage(chatId, messages.UnknownCommand)
		return
	}

	unauthorizedUserId, err := b.db.GetAnyUnauthorizedUser()
	if err != nil {
		log.Println("Error fetching row from database:", err)
		b.sendMessage(chatId, messages.GetUnathorizedUsersError)
		return
	}

	if unauthorizedUserId == -1 {
		b.sendMessage(chatId, messages.AllUsersAuthorized)
		return
	}

	b.sendMessageWithButtons(chatId, unauthorizedUserId)
}
