package tg

import (
	"log"
	"taskbot/delivery/tg/messages"

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

func (r *Responder) Send(m messages.Message) {
	msgConf := tgbotapi.NewMessage(m.ChatId, m.Text)
	msgConf.ParseMode = tgbotapi.ModeHTML

	_, err := r.bot.Send(msgConf)

	if err != nil {
		log.Printf("bot send error. chat: %v, err: %v", m.ChatId, err)
	}
}
