package messages

import (
	"taskbot/domain"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func HelloMessage(user domain.User) tgbotapi.MessageConfig {
	text := `
Создай задачу. Выполни задачу. Съешь печенку.
Создать задачу: /create
	`

	return tgbotapi.NewMessage(user.TelegramId, text)
}
