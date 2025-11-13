package tg

import (
	"log"
	"taskbot/service/telegram"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type UpdateHandler struct {
	resp      *Responder
	processor *telegram.UpdateProcessor
}

func NewUpdateHandler(r *Responder, processor *telegram.UpdateProcessor) *UpdateHandler {
	return &UpdateHandler{
		resp:      r,
		processor: processor,
	}
}

func (h *UpdateHandler) Handle(update tgbotapi.Update) {

	h.logUpdate(update)
	msges, err := h.processor.Handle(update)
	if err != nil {
		log.Printf("error processing update: %v\n", err)
		return
	}

	for _, m := range msges {
		err = h.resp.Send(m)
		if err != nil {
			log.Printf("bot send error. chat: %v, err: %v\n", m.ChatID, err)
		}
	}
}

func (h *UpdateHandler) logUpdate(update tgbotapi.Update) {
	chatId := update.FromChat().ID

	// var data string
	// if update.Message != nil {
	// 	data = update.Message.Text
	// } else if update.CallbackQuery != nil {
	// 	data = update.CallbackData()
	// }
	// log.Printf("got update. chat: %v. data: %v", chatId, data)

	log.Printf("got update. chat: %v. update: %v", chatId, update)
}
