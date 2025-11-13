package telegram

import (
	"errors"
	"fmt"
	"taskbot/domain"
	"taskbot/repository"
	"taskbot/service"
	"taskbot/service/task"
	"taskbot/service/telegram/keyboard"
	"taskbot/service/telegram/messages"
	"taskbot/service/user"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	tasksInPage = 10
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

func (p *UpdateProcessor) Handle(update tgbotapi.Update) ([]tgbotapi.MessageConfig, error) {
	u, err := p.extractUser(update)
	if err != nil {
		return []tgbotapi.MessageConfig{}, err
	}

	state, err := p.getChatState(u)
	if err != nil {
		return []tgbotapi.MessageConfig{}, err
	}

	processUpdate, err := p.extractProcessUpdate(update)
	if err != nil {
		return []tgbotapi.MessageConfig{}, err
	}

	return p.process(u, state, processUpdate)
}

func (p *UpdateProcessor) process(user domain.User, state ChatState, update Update) ([]tgbotapi.MessageConfig, error) {
	updateData := update.GetData()
	if update.isCommand {
		return p.handleCommand(user, state, updateData)
	}

	if update.isCallback {
		return p.handleCallback(user, state, updateData)
	}

	return p.handleTextMessage(user, state, updateData)
}

func (p *UpdateProcessor) handleCommand(user domain.User, state ChatState, command string) ([]tgbotapi.MessageConfig, error) {
	msges := []tgbotapi.MessageConfig{}

	switch command {
	case "start", "/start":
		msg := tgbotapi.NewMessage(user.TelegramId, messages.HelloMessage())
		msges = append(msges, msg)
	case "create", "/create":
		msg := tgbotapi.NewMessage(user.TelegramId, messages.WaitTaskTitle())
		task := domain.Task{
			UserId: user.Id,
			Status: domain.WAITING,
		}

		state.Data = task
		state.Status = WAIT_TASK_TITLE
		err := p.saveChatState(state)
		if err != nil {
			return msges, err
		}

		msges = append(msges, msg)
	case "list", "/list":
		tasks, err := p.getTasksPerPage(user, 1)
		if err != nil {
			return msges, err
		}

		msges = p.getTasksMessages(user, tasks)
	}

	return msges, nil
}

func (p *UpdateProcessor) handleCallback(user domain.User, state ChatState, callback string) ([]tgbotapi.MessageConfig, error) {
	// parse callback -> tasks callback ->
	// next page callback + page -> send list task for next page
	// ready task callback + task id -> update task -> ready task message
	// delete task callback + task id -> delete task -> deleted task message

	msgConf := []tgbotapi.MessageConfig{}
	return msgConf, nil
}

func (p *UpdateProcessor) handleTextMessage(user domain.User, state ChatState, text string) ([]tgbotapi.MessageConfig, error) {
	// check state status ->
	// update state data for step + "next waiting data" message or "task message" + callback

	msgConf := []tgbotapi.MessageConfig{}
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
		err = p.saveChatState(state)
	}

	return state, err
}

func (p *UpdateProcessor) saveChatState(state ChatState) error {
	state.UpdateAt = time.Now()
	return p.stateStore.Save(state)
}

func (p *UpdateProcessor) getTasksPerPage(user domain.User, page int) ([]domain.Task, error) {
	if page < 1 {
		page = 1
	}

	var pagedTasks []domain.Task

	tasks, err := p.taskService.GetAll(user)
	if err != nil {
		return pagedTasks, err
	}

	pagedTasks = make([]domain.Task, 0, 10)

	start := (page - 1) * tasksInPage
	end := page * tasksInPage

	if start >= len(tasks) {
		return []domain.Task{}, nil
	}

	if end > len(tasks) {
		end = len(tasks)
	}

	return pagedTasks[start:end], nil
}

func (p *UpdateProcessor) getTasksMessages(user domain.User, tasks []domain.Task) []tgbotapi.MessageConfig {
	msges := make([]tgbotapi.MessageConfig, 0, len(tasks))

	for _, t := range tasks {
		msg := tgbotapi.NewMessage(user.TelegramId, messages.TaskContent(t))
		msg.ReplyMarkup = keyboard.TaskInlineKeyboard(t)

		msges = append(msges, msg)
	}

	return msges
}
