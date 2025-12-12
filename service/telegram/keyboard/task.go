package keyboard

import (
	"strconv"
	"taskbot/domain"
	"taskbot/service/telegram/callback"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func TaskInlineKeyboard(task domain.Task) tgbotapi.InlineKeyboardMarkup {
	keyboard := tgbotapi.InlineKeyboardMarkup{}

	strId := strconv.FormatInt(task.Id, 10)
	rowBtns := []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData("Выполнена", callback.CreateCallbackData(callback.TaskDone, strId)),
		tgbotapi.NewInlineKeyboardButtonData("Удалить", callback.CreateCallbackData(callback.TaskDelete, strId)),
	}

	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, rowBtns)

	return keyboard
}

func NextPageInlineKeyboard(text string, page int) tgbotapi.InlineKeyboardMarkup {
	keyboard := tgbotapi.InlineKeyboardMarkup{}
	strPage := strconv.Itoa(page)
	rowBtns := []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData(text, callback.CreateCallbackData(callback.NextTasksPage, strPage)),
	}

	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, rowBtns)

	return keyboard
}
