package tg

import (
	"log"
	"taskbot/delivery/tg/messages"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type UpdateHandler struct {
	resp *Responder
}

func NewUpdateHandler(r *Responder) *UpdateHandler {
	return &UpdateHandler{
		resp: r,
	}
}

func (h *UpdateHandler) Handle(update tgbotapi.Update) {
	if update.Message == nil {
		log.Printf("skip empty message, chat: %v", update.Message.Chat.ID)
		return
	}

	if update.Message.IsCommand() {
		log.Printf("get command. chat_id: %v, command: %.8v", update.Message.Chat.ID, update.Message.Command())

		command := update.Message.Command()
		chatId := update.Message.Chat.ID

		switch command {
		case "start":
			h.resp.Send(messages.Message{
				ChatId: chatId,
				Text:   messages.HelloMessageText(),
			})
		}
		//command handling
	} else {
		//inline text handling
	}

}
