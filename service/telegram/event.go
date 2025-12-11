package telegram

import "fmt"

type Event struct {
	ChatID int64
	Data   string
	Type   EventType
}

type EventType int

const (
	EventTypeUnknown EventType = iota
	EventTypeCommand
	EventTypeCallback
	EventTypeText
)

func NewCommandEvent(chatID int64, command string) Event {
	return Event{
		ChatID: chatID,
		Data:   command,
		Type:   EventTypeCommand,
	}
}

func NewCallbackEvent(chatID int64, callbackData string) Event {
	return Event{
		ChatID: chatID,
		Data:   callbackData,
		Type:   EventTypeCallback,
	}
}

func NewTextEvent(chatID int64, text string) Event {
	return Event{
		ChatID: chatID,
		Data:   text,
		Type:   EventTypeText,
	}
}

func (e Event) IsCommand() bool {
	return e.Type == EventTypeCommand
}

func (e Event) IsCallback() bool {
	return e.Type == EventTypeCallback
}

func (e Event) IsText() bool {
	return e.Type == EventTypeText
}

func (e Event) IsValid() bool {
	return e.ChatID != 0 && e.Data != "" && e.Type != EventTypeUnknown
}

func (e Event) String() string {
	typeNames := map[EventType]string{
		EventTypeCommand:  "Command",
		EventTypeCallback: "Callback",
		EventTypeText:     "Text",
		EventTypeUnknown:  "Unknown",
	}
	return fmt.Sprintf("%s{chat:%d, data:%q}", typeNames[e.Type], e.ChatID, e.Data)
}
