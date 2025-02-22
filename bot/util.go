package bot

import (
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/yalexaner/nag-meetings-go/database"
	"github.com/yalexaner/nag-meetings-go/messages"
)

func (b *Bot) sendMessage(chatId int64, text string) {
	msg := tgbotapi.NewMessage(chatId, text)
	_, err := b.api.Send(msg)
	if err != nil {
		log.Printf("Error sending message: %v", err)
	}
}

func (b *Bot) editMessage(chatId int64, messageId int, text string) {
	editMsg := tgbotapi.NewEditMessageText(chatId, messageId, text)
	_, err := b.api.Send(editMsg)
	if err != nil {
		log.Printf("Error editing message: %v", err)
	}
}

func (b *Bot) sendMessageWithButtons(chatId int64, user *database.User) {
	if user == nil {
		log.Println("User to be send is nil")
		return
	}

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(messages.Yes, fmt.Sprintf("%d_1", user.UserID)),
			tgbotapi.NewInlineKeyboardButtonData(messages.No, fmt.Sprintf("%d_0", user.UserID)),
		),
	)

	msg := tgbotapi.NewMessage(chatId, fmt.Sprintf(messages.AuthorizeUserQuestion, user.Name))
	msg.ReplyMarkup = keyboard

	_, err := b.api.Send(msg)
	if err != nil {
		log.Println("Error sending message with buttons:", err)
	}
}

func (b *Bot) editMessageWithButtons(chatId int64, messageId int, user *database.User) {
	if user == nil {
		log.Println("User to be send is nil")
		return
	}

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(messages.Yes, fmt.Sprintf("%d_1", user.UserID)),
			tgbotapi.NewInlineKeyboardButtonData(messages.No, fmt.Sprintf("%d_0", user.UserID)),
		),
	)

	editMsg := tgbotapi.NewEditMessageText(chatId, messageId, fmt.Sprintf(messages.AuthorizeUserQuestion, user.Name))
	editMsg.ReplyMarkup = &keyboard

	_, err := b.api.Send(editMsg)
	if err != nil {
		log.Println("Error editing message with buttons:", err)
	}
}
