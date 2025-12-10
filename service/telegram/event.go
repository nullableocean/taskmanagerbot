package telegram

type Event struct {
	chatId int64
	data   string

	isCommand, isCallback, isText bool
}

func (u Event) GetChatId() int64 {
	return u.chatId
}

func (u Event) GetData() string {
	return u.data
}

func (u Event) IsText() bool {
	return u.isText
}

func (u Event) IsCommand() bool {
	return u.isCommand
}

func (u Event) IsCallback() bool {
	return u.isCallback
}

func NewCallbackEvent(chatId int64, callback string) Event {
	return Event{
		chatId:     chatId,
		data:       callback,
		isCommand:  false,
		isCallback: true,
		isText:     false,
	}
}

func NewCommandEvent(chatId int64, command string) Event {
	return Event{
		chatId:     chatId,
		data:       command,
		isCommand:  false,
		isCallback: true,
		isText:     false,
	}
}
func NewTextEvent(chatId int64, text string) Event {
	return Event{
		chatId:     chatId,
		data:       text,
		isCommand:  false,
		isCallback: true,
		isText:     false,
	}
}
