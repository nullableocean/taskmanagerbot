package domain

import "time"

type TaskStatus string

const (
	WAITING TaskStatus = "WAIT"
	READY   TaskStatus = "READY"
)

type Task struct {
	Id     int64      `json:"id,omitemtpy"`
	Title  string     `json:"title"`
	Body   string     `json:"body"`
	Status TaskStatus `json:"status"`

	UserId    int64     `json:"user_id"`
	CreatedAt time.Time `json:"created_at,omitemtpy"`
	UpdatedAt time.Time `json:"updated_at,omitemtpy"`
}
