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
		msg := messages.TaskContent(user, t)
		msg.ReplyMarkup = keyboard.TaskInlineKeyboard(t)

		msges = append(msges, msg)
	}

	return msges
}
