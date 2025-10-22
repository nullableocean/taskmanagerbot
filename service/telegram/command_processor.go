package telegram

import (
	"taskbot/repository"
	"taskbot/service/user"
)

type CommandProcessor struct {
	utgSer *user.TelegramUserService
	srepo  repository.StateRepository
}
