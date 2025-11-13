package telegram

type Update struct {
	chatId int64
	data   string

	isCommand, isCallback, isText bool
}

func (u Update) GetChatId() int64 {
	return u.chatId
}

func (u Update) GetData() string {
	return u.data
}

func (u Update) IsText() bool {
	return u.isText
}

func (u Update) IsCommand() bool {
	return u.isCommand
}

func (u Update) IsCallback() bool {
	return u.isCallback
}

func NewCallbackUpdate(chatId int64, callback string) Update {
	return Update{
		chatId:     chatId,
		data:       callback,
		isCommand:  false,
		isCallback: true,
		isText:     false,
	}
}

func NewCommandUpdate(chatId int64, command string) Update {
	return Update{
		chatId:     chatId,
		data:       command,
		isCommand:  false,
		isCallback: true,
		isText:     false,
	}
}
func NewTextUpdate(chatId int64, text string) Update {
	return Update{
		chatId:     chatId,
		data:       text,
		isCommand:  false,
		isCallback: true,
		isText:     false,
	}
}
