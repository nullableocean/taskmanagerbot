package telegram

import (
	"errors"
	"fmt"
	"taskbot/domain"
	"taskbot/repository"
	"taskbot/service"
	"taskbot/service/task"
	"taskbot/service/user"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type StateStore interface {
	Get(id int64) (ChatState, error)
	Save(ChatState) error
	Delete(ChatState) error
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

func (p *UpdateProcessor) Handle(update tgbotapi.Update) (tgbotapi.MessageConfig, error) {
	u, err := p.extractUser(update)
	if err != nil {
		return tgbotapi.MessageConfig{}, err
	}

	state, err := p.getChatState(u)
	if err != nil {
		return tgbotapi.MessageConfig{}, err
	}

	processUpdate, err := p.extractProcessUpdate(update)
	if err != nil {
		return tgbotapi.MessageConfig{}, err
	}

	return p.process(u, state, processUpdate)
}

func (p *UpdateProcessor) process(user domain.User, state ChatState, update Update) (tgbotapi.MessageConfig, error) {
	//1) check state -> handling
	//2) check updateType; callback -> handle callback(user, state, callback); command -> handle command(user, state, command)
	// text -> handleText(user, state, text)
	updateData := update.GetData()
	if update.isCommand {
		return p.handleCommand(user, state, updateData)
	}

	if update.isCallback {
		return p.handleCallback(user, state, updateData)
	}

	return p.handleTextMessage(user, state, updateData)
}

func (p *UpdateProcessor) handleCommand(user domain.User, state ChatState, command string) (tgbotapi.MessageConfig, error) {
	switch command {
	case "start", "/start":
		// info message
	case "create", "/create":
		// create task
		// set state for creating -> message "create task" -> continue handle in text handler
	case "list", "/list":
		// create task list, offset page 10
		// create callbacks for tasks, callback next page, set state listed tasks
		// send list message or array messages "task info + callback"?
	}

	msgConf := tgbotapi.MessageConfig{}
	return msgConf, nil
}

func (p *UpdateProcessor) handleTextMessage(user domain.User, state ChatState, text string) (tgbotapi.MessageConfig, error) {
	// check state status ->
	// update state data for step + "next waiting data" message or "task message" + callback

	msgConf := tgbotapi.MessageConfig{}
	return msgConf, nil
}

func (p *UpdateProcessor) handleCallback(user domain.User, state ChatState, callback string) (tgbotapi.MessageConfig, error) {
	// parse callback -> tasks callback ->
	// next page callback + page -> send list task for next page
	// ready task callback + task id -> update task -> ready task message
	// delete task callback + task id -> delete task -> deleted task message

	msgConf := tgbotapi.MessageConfig{}
	return msgConf, nil
}

func (p *UpdateProcessor) extractUser(update tgbotapi.Update) (domain.User, error) {
	var chatId int64

	user, err := p.userTgService.FindByTelegramId(chatId)
	if err != nil && errors.As(err, service.ErrNotFound) {
		return p.userTgService.CreateFromUpdate(update)
	}

	return user, err
}

func (p *UpdateProcessor) extractProcessUpdate(update tgbotapi.Update) (Update, error) {
	var processUpdate Update

	switch {
	case update.CallbackQuery != nil:
		chatId := update.CallbackQuery.From.ID
		processUpdate = NewCallbackUpdate(chatId, update.CallbackData())
	case update.Message != nil:
		chatId := update.Message.Chat.ID

		if update.Message.IsCommand() {
			processUpdate = NewCommandUpdate(chatId, update.Message.Command())
		} else {
			processUpdate = NewTextUpdate(chatId, update.Message.Text)
		}
	default:
		return processUpdate, fmt.Errorf("unknown update. chat: %v, data: %v\n", update.FromChat().ID, update)
	}

	return processUpdate, nil
}

func (p *UpdateProcessor) getChatState(user domain.User) (ChatState, error) {
	state, err := p.stateStore.Get(user.TelegramId)
	if errors.As(err, repository.ErrNotFound) {
		state = ChatState{
			Id:       user.TelegramId,
			Status:   IDLE,
			UpdateAt: time.Now(),
			Data:     nil,
		}
		err = p.setChatState(state)
	}

	return state, err
}

func (p *UpdateProcessor) setChatState(state ChatState) error {
	return p.stateStore.Save(state)
}
