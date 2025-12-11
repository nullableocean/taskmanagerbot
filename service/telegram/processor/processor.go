package processor

import (
	"errors"
	"log"
	"taskbot/domain"
	"taskbot/repository"
	"taskbot/service/task"
	"taskbot/service/telegram"
	"taskbot/service/user"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	tasksInPage = 10
)

type StateStore interface {
	Get(id int64) (telegram.ChatState, error)
	Save(telegram.ChatState) error
	Delete(telegram.ChatState) error
}

type UpdateProcessor struct {
	userTgService *user.TelegramUserService
	taskService   *task.TaskService
	stateStore    StateStore
}

func NewUpdateProccesor(userTgService *user.TelegramUserService, tService *task.TaskService, stateStore StateStore) *UpdateProcessor {
	return &UpdateProcessor{
		userTgService: userTgService,
		taskService:   tService,
		stateStore:    stateStore,
	}
}

func (p *UpdateProcessor) Handle(update tgbotapi.Update) ([]tgbotapi.MessageConfig, error) {
	u, err := p.extractUserFromUpdate(update)
	if err != nil {
		return []tgbotapi.MessageConfig{}, err
	}

	state, err := p.getChatState(u)
	if err != nil {
		return []tgbotapi.MessageConfig{}, err
	}

	event, err := p.extractEventFromUpdate(update)
	if err != nil {
		return []tgbotapi.MessageConfig{}, err
	}

	return p.process(u, state, event)
}

func (p *UpdateProcessor) process(user domain.User, state telegram.ChatState, event telegram.Event) ([]tgbotapi.MessageConfig, error) {
	eventData := event.Data
	if event.IsCommand() {
		return p.handleCommand(user, state, eventData)
	}

	if event.IsCallback() {
		return p.handleCallback(user, state, eventData)
	}

	return p.handleTextMessage(user, state, eventData)
}

func (p *UpdateProcessor) getChatState(user domain.User) (telegram.ChatState, error) {
	state, err := p.stateStore.Get(user.TelegramId)
	if errors.Is(err, repository.ErrNotFound) {
		state = telegram.ChatState{
			Id:       user.TelegramId,
			Status:   telegram.IDLE,
			UpdateAt: time.Now(),
			Data:     nil,
		}
		err = p.saveChatState(state)
	}

	log.Println("state getted", state.Status, state.Id, state.Data)

	return state, err
}

func (p *UpdateProcessor) saveChatState(state telegram.ChatState) error {
	state.UpdateAt = time.Now()
	return p.stateStore.Save(state)
}
