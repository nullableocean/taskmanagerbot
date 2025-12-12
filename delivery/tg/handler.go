package tg

import (
	"log"
	"runtime/debug"
	"strings"
	"taskbot/service/telegram/processor"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type UpdateHandler struct {
	resp      *Responder
	processor *processor.UpdateProcessor
}

func NewUpdateHandler(r *Responder, processor *processor.UpdateProcessor) *UpdateHandler {
	return &UpdateHandler{
		resp:      r,
		processor: processor,
	}
}

func (h *UpdateHandler) Handle(update tgbotapi.Update) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("recovered panic: ", r)

			stack := debug.Stack()

			lines := strings.Split(string(stack), "\n")
			for i, line := range lines {
				log.Printf("[%d] %s", i, line)
			}

			h.sendErrorMessage(update.FromChat().ID)
		}
	}()

	if update.FromChat() == nil {
		log.Println("update chat is nil")
		return
	}
	h.logUpdate(update)

	msges, err := h.processor.Handle(update)
	if err != nil {
		log.Printf("error processing update: %v\n", err)
		h.sendErrorMessage(update.FromChat().ID)

		return
	}

	for _, m := range msges {
		err = h.resp.Send(m)
		if err != nil {
			log.Printf("bot send error. chat: %v, err: %v\n", m.ChatID, err)
		}
	}
}

func (h *UpdateHandler) sendErrorMessage(chatId int64) {
	h.resp.Send(tgbotapi.NewMessage(chatId, "Что-то пошло не так. Попробуй ещё раз :)"))
}

func (h *UpdateHandler) logUpdate(update tgbotapi.Update) {
	log.Printf("got update. user: %v", update.FromChat().UserName)
}
