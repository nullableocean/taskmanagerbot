package tg

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Responder struct {
	bot *tgbotapi.BotAPI
}

func NewResponder(bot *tgbotapi.BotAPI) *Responder {
	return &Responder{
		bot: bot,
	}
}

func (r *Responder) Send(msgConfig tgbotapi.MessageConfig) error {
	msgConfig.ParseMode = tgbotapi.ModeHTML

	_, err := r.bot.Send(msgConfig)

	return err
}
