package processor

import (
	"taskbot/domain"
	"taskbot/service/telegram/keyboard"
	"taskbot/service/telegram/messages"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (p *UpdateProcessor) getTasksPerPage(user domain.User, page int) ([]domain.Task, error) {
	if page < 1 {
		page = 1
	}

	tasks, err := p.taskService.GetAll(user)
	if err != nil {
		return nil, err
	}

	start := (page - 1) * tasksInPage
	end := page * tasksInPage

	if start >= len(tasks) {
		return []domain.Task{}, nil
	}

	if end > len(tasks) {
		end = len(tasks)
	}

	return tasks[start:end], nil
}

func (p *UpdateProcessor) getTasksMessages(user domain.User, tasks []domain.Task) []tgbotapi.MessageConfig {
	msges := make([]tgbotapi.MessageConfig, 0, len(tasks))

	for _, t := range tasks {
		msg := messages.TaskContent(user, t)
		msg.ReplyMarkup = keyboard.TaskInlineKeyboard(t)

		msges = append(msges, msg)
	}

	return msges
}

func (p *UpdateProcessor) getTasksListMessages(user domain.User, page int) ([]tgbotapi.MessageConfig, error) {
	tasks, err := p.getTasksPerPage(user, page)
	if err != nil {
		return nil, err
	}

	msges := p.getTasksMessages(user, tasks)

	nextPageMsg := tgbotapi.NewMessage(user.TelegramId, "")
	if len(tasks) != 0 {
		nextPageMsg.Text = messages.NextPageMessage
		nextPageMsg.ReplyMarkup = keyboard.NextPageInlineKeyboard(messages.NextBtnMessage, page+1)
	} else {
		backPage := page - 1
		if backPage < 1 {
			backPage = 1
		}

		nextPageMsg.Text = messages.NextPageEmptyMessage
		nextPageMsg.ReplyMarkup = keyboard.NextPageInlineKeyboard(messages.BackBtnMessage, backPage)
	}

	msges = append(msges, nextPageMsg)

	return msges, nil
}
