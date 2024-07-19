package bot

import (
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
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

func (b *Bot) sendMessageWithButtons(chatId int64, unauthorizedUserId int64) {
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(messages.Yes, fmt.Sprintf("%d_1", unauthorizedUserId)),
			tgbotapi.NewInlineKeyboardButtonData(messages.No, fmt.Sprintf("%d_0", unauthorizedUserId)),
		),
	)

	msg := tgbotapi.NewMessage(chatId, messages.AuthorizeUserQuestion)
	msg.ReplyMarkup = keyboard

	_, err := b.api.Send(msg)
	if err != nil {
		log.Println("Error sending message with buttons:", err)
	}
}

func (b *Bot) editMessageWithButtons(chatId int64, messageId int, unauthorizedUserId int64) {
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(messages.Yes, fmt.Sprintf("%d_1", unauthorizedUserId)),
			tgbotapi.NewInlineKeyboardButtonData(messages.No, fmt.Sprintf("%d_0", unauthorizedUserId)),
		),
	)

	editMsg := tgbotapi.NewEditMessageText(chatId, messageId, messages.AuthorizeUserQuestion)
	editMsg.ReplyMarkup = &keyboard

	_, err := b.api.Send(editMsg)
	if err != nil {
		log.Println("Error editing message with buttons:", err)
	}
}
