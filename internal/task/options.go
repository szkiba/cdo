package task

import (
	"strings"

	"github.com/google/shlex"
)

func getopts(task *Task, opts map[string]string) {
	for key, value := range opts {
		switch strings.ToLower(key) {
		case "requires":
			parts := strings.Split(value, ",")
			for _, part := range parts {
				args, err := shlex.Split(part)
				if err == nil {
					task.Requires = append(task.Requires, args)
				}
			}
		default:
		}
	}
}
