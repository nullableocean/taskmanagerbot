package processor

import (
	"errors"
	"strconv"
	"taskbot/domain"
	"taskbot/service/telegram"
	"taskbot/service/telegram/callback"
	"taskbot/service/telegram/messages"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (p *UpdateProcessor) handleTextMessage(user domain.User, state telegram.ChatState, text string) ([]tgbotapi.MessageConfig, error) {
	msges := []tgbotapi.MessageConfig{}

	switch state.Status {
	case telegram.WAIT_TASK_TITLE:
		t, ok := state.Data.(domain.Task)
		if !ok {
			return nil, errors.New("invalid data in state")
		}

		t.Title = text
		state.Data = t
		state.Status = telegram.WAIT_TASK_BODY

		err := p.stateStore.Save(state)
		if err != nil {
			return nil, err
		}

		msges = append(msges, messages.WaitTaskBody(user))
	case telegram.WAIT_TASK_BODY:
		t, ok := state.Data.(domain.Task)
		if !ok {
			return nil, errors.New("invalid data in state")
		}

		t.Body = text

		t, err := p.taskService.Create(user, t)
		if err != nil {
			return nil, err
		}

		state.Status = telegram.IDLE
		err = p.stateStore.Save(state)
		if err != nil {
			return nil, err
		}

		msges = append(msges, messages.TaskCreated(user), messages.TaskContent(user, t))
	default:
		msges = append(msges, messages.HelloMessage(user))
	}

	return msges, nil
}

func (p *UpdateProcessor) handleCallback(user domain.User, state telegram.ChatState, callbackData string) ([]tgbotapi.MessageConfig, error) {
	msges := []tgbotapi.MessageConfig{}

	op, data := callback.ExtractCallbackData(callbackData)
	switch op {
	case callback.TaskDone:
		taskId, err := strconv.ParseInt(data, 10, 64)
		if err != nil {
			return nil, err
		}

		t, err := p.taskService.GetById(taskId)
		if err != nil {
			return nil, err
		}

		var msg tgbotapi.MessageConfig
		if t.Status != domain.READY {
			t.Status = domain.READY
			_, err = p.taskService.Update(t)

			if err != nil {
				return nil, err
			}

			msg = messages.TaskReady(user)
		} else {
			msg = messages.TaskAlreadyReady(user)
		}

		msges = append(msges, msg)

	case callback.TaskDelete:
		taskId, err := strconv.ParseInt(data, 10, 64)
		if err != nil {
			return nil, err
		}

		t, err := p.taskService.GetById(taskId)
		if err != nil {
			return nil, err
		}

		err = p.taskService.Delete(t)
		if err != nil {
			return nil, err
		}

		msges = append(msges, messages.TaskDeleted(user))

	case callback.NextTasksPage:
		nextPage, err := strconv.Atoi(data)
		if err != nil {
			return nil, err
		}

		msges, err = p.getTasksListMessages(user, nextPage)
		if err != nil {
			return nil, err
		}
	}

	return msges, nil
}

func (p *UpdateProcessor) handleCommand(user domain.User, state telegram.ChatState, command string) ([]tgbotapi.MessageConfig, error) {
	msges := []tgbotapi.MessageConfig{}

	switch command {
	case "start", "/start":
		msges = append(msges, messages.HelloMessage(user))
	case "create", "/create":
		task := domain.Task{
			UserId: user.Id,
			Status: domain.WAITING,
		}

		state.Data = task
		state.Status = telegram.WAIT_TASK_TITLE
		err := p.saveChatState(state)
		if err != nil {
			return nil, err
		}

		msges = append(msges, messages.WaitTaskBody(user))
	case "list", "/list":
		var err error

		msges, err = p.getTasksListMessages(user, 1)
		if err != nil {
			return nil, err
		}
	}

	return msges, nil
}
