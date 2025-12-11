package telegram

import (
	"encoding/json"
	"taskbot/domain"
	"time"
)

type Status string

const (
	IDLE            Status = "idle"
	WAIT_TASK_TITLE Status = "wait_title"
	WAIT_TASK_BODY  Status = "wait_body"
)

type ChatState struct {
	Id       int64           `json:"id"`
	Status   Status          `json:"status"`
	UpdateAt time.Time       `json:"update_at"`
	Data     json.RawMessage `json:"data,omitempty"`
}

func (cs ChatState) GetTask() (*domain.Task, error) {
	if len(cs.Data) == 0 {
		return nil, nil
	}

	var task domain.Task
	if err := json.Unmarshal(cs.Data, &task); err != nil {
		return nil, err
	}
	return &task, nil
}

func (cs *ChatState) SetTask(task domain.Task) error {
	data, err := json.Marshal(task)
	if err != nil {
		return err
	}
	cs.Data = data
	return nil
}

func (cs *ChatState) ClearData() {
	cs.Data = nil
}
