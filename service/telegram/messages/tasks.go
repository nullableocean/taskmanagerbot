package messages

import (
	"fmt"
	"taskbot/domain"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func TaskContent(user domain.User, task domain.Task) tgbotapi.MessageConfig {
	format := `
<b>%s</b>

%s

<i>%s</i>
`

	var status string
	switch task.Status {
	case domain.READY:
		status = "Выполнена"
	default:
		status = "Ожидает"
	}

	text := fmt.Sprintf(format, task.Title, task.Body, status)

	return tgbotapi.NewMessage(user.TelegramId, text)
}

func WaitTaskTitle(user domain.User) tgbotapi.MessageConfig {
	text := `
<b>Озаглавь задачу:</b>	
`

	return tgbotapi.NewMessage(user.TelegramId, text)
}

func WaitTaskBody(user domain.User) tgbotapi.MessageConfig {
	text := `
<b>Опиши суть задачи:</b>	
`

	return tgbotapi.NewMessage(user.TelegramId, text)
}

func TaskCreated(user domain.User) tgbotapi.MessageConfig {
	text := `
<b>Задача успешно создана.</b>	
`

	return tgbotapi.NewMessage(user.TelegramId, text)
}

func TaskReady(user domain.User) tgbotapi.MessageConfig {
	text := `
<b>Задача выполнена.</b>	
`

	return tgbotapi.NewMessage(user.TelegramId, text)
}

func TaskAlreadyReady(user domain.User) tgbotapi.MessageConfig {
	text := `
<b>Задача уже помечена, как выполненная.</b>	
`

	return tgbotapi.NewMessage(user.TelegramId, text)
}
