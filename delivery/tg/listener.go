package tg

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type ListenerState int

const (
	WAIT ListenerState = iota
	RUNNED
	STOPED
)

type UpdateListener struct {
	updatesCh tgbotapi.UpdatesChannel
	state     ListenerState
	stopCH    chan struct{}

	handler *UpdateHandler
}

func NewUpdateListener(updates tgbotapi.UpdatesChannel, h *UpdateHandler) *UpdateListener {
	return &UpdateListener{
		updatesCh: updates,
		state:     WAIT,
		stopCH:    make(chan struct{}),
		handler:   h,
	}
}

func (l *UpdateListener) Listen() error {
	if l.state == RUNNED {
		return ErrListenerRunned
	}

LOOP:
	for {
		select {
		case u := <-l.updatesCh:
			go l.handler.Handle(u)
		case <-l.stopCH:
			log.Println("listener stoped")

			l.state = STOPED
			break LOOP
		}
	}

	return nil
}

func (l *UpdateListener) Stop() {
	if l.state == RUNNED {
		l.stopCH <- struct{}{}
	}
}
