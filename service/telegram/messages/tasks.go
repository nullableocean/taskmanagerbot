package messages

import (
	"fmt"
	"taskbot/domain"
)

const ()

func WaitTaskTitle() string {
	return `
<b>Озаглавь задачу:</b>	
`
}

func WaitTaskBody() string {
	return `
<b>Опиши суть задачи:</b>	
`
}

func TaskContent(task domain.Task) string {
	format := `
<b>%s</b>

%s

<i>%s</i>
`

	var status string
	switch task.Status {
	case domain.READY:
		status = "Выполнена"
	default:
		status = "Ожидает"
	}

	return fmt.Sprintf(format, task.Title, task.Body, status)
}
