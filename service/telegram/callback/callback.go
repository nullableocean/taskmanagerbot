package callback

import (
	"fmt"
	"strings"
)

type Operation string

const (
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
	return Operation(parts[0]), parts[1]
}
