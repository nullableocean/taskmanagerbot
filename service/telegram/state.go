package telegram

import (
	"time"
)

type Status string

const (
	IDLE            Status = "idle"
	WAIT_TASK_TITLE Status = "wait_title"
	WAIT_TASK_BODY  Status = "wait_body"
)

type ChatState struct {
	Id       int64       `json:"id"`
	Status   Status      `json:"status"`
	UpdateAt time.Time   `json:"update_at"`
	Data     interface{} `json:"data"`
}
