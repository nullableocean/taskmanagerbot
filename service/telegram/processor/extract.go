package processor

import (
	"errors"
	"fmt"
	"taskbot/domain"
	"taskbot/service"
	"taskbot/service/telegram"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (p *UpdateProcessor) extractUserFromUpdate(update tgbotapi.Update) (domain.User, error) {
	var chatId int64

	user, err := p.userTgService.FindByTelegramId(chatId)
	if err != nil && errors.As(err, service.ErrNotFound) {
		return p.userTgService.CreateFromUpdate(update)
	}

	return user, err
}

func (p *UpdateProcessor) extractEventFromUpdate(update tgbotapi.Update) (telegram.Event, error) {
	var processUpdate telegram.Event

	switch {
	case update.CallbackQuery != nil:
		chatId := update.CallbackQuery.From.ID
		processUpdate = telegram.NewCallbackEvent(chatId, update.CallbackData())
	case update.Message != nil:
		chatId := update.Message.Chat.ID

		if update.Message.IsCommand() {
			processUpdate = telegram.NewCommandEvent(chatId, update.Message.Command())
		} else {
			processUpdate = telegram.NewTextEvent(chatId, update.Message.Text)
		}
	default:
		return processUpdate, fmt.Errorf("unknown update. chat: %v, data: %v\n", update.FromChat().ID, update)
	}

	return processUpdate, nil
}
