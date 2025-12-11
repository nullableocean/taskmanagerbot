package callback

import (
	"fmt"
	"strings"
)

type Operation string

const (
	Undefined     Operation = "undefined"
	TaskDone      Operation = "tdone"
	TaskDelete    Operation = "tdelete"
	NextTasksPage Operation = "tnextpage"
)

const (
	separate = ":"
)

func CreateCallbackData(op Operation, data string) string {
	return fmt.Sprintf("%s%s%s", op, separate, data)
}

func ExtractCallbackData(callback string) (Operation, string) {
	parts := strings.Split(callback, separate)

	if len(parts) > 1 {
		return Operation(parts[0]), parts[1]
	}

	return Undefined, ""
}
