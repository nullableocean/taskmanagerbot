package repository

import "taskbot/domain"

type TelegramUserRepository interface {
	GetByTelegramId(tid int64) (domain.User, error)
}

type UserRepository interface {
	Save(domain.User) (domain.User, error)
	Get(id int64) (domain.User, error)
}
