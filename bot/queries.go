package bot

import (
	"log"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/yalexaner/nag-meetings-go/messages"
)

func (b *Bot) handleCallbackQuery(query *tgbotapi.CallbackQuery) {
	parts := strings.Split(query.Data, "_")
	if len(parts) != 2 {
		log.Println("Invalid callback data format")
		return
	}

	id, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		log.Println("Error parsing ID:", err)
		return
	}

	var value bool
	if parts[1] == "1" {
		value = true
	} else {
		value = false
	}

	err = b.db.ChangeAuthorization(id, value)
	if err != nil {
		log.Println("Error processing row data:", err)
		callback := tgbotapi.NewCallbackWithAlert(query.ID, messages.ChangeAuthorizationErrror)
		if _, err := b.api.AnswerCallbackQuery(callback); err != nil {
			log.Println("Error answering callback query:", err)
		}
		return
	}

	unauthorizedUserId, err := b.db.GetAnyUnauthorizedUser()
	if err != nil {
		log.Println("Error fetching row from database:", err)
		b.sendMessage(query.Message.Chat.ID, messages.GetUnathorizedUsersError)
		return
	}

	if unauthorizedUserId == -1 {
		b.editMessage(query.Message.Chat.ID, query.Message.MessageID, messages.AllUsersAuthorized)
		return
	}

	b.editMessageWithButtons(query.Message.Chat.ID, query.Message.MessageID, unauthorizedUserId)

	callback := tgbotapi.NewCallback(query.ID, "")
	if _, err := b.api.AnswerCallbackQuery(callback); err != nil {
		log.Println("Error answering callback query:", err)
	}
}
