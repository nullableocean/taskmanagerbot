package domain

import "time"

type User struct {
	Id       int64
	Username string
	Name     string
	Password string

	TelegramId int64

	CreatedAt time.Time
	UpdatedAt time.Time
}
