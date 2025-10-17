package repository

import "taskbot/domain"

type UserRepository interface {
	Create(domain.User) error
	Get(id int64) (domain.User, error)
	GetByTelegramId(tid int64) (domain.User, error)
}
