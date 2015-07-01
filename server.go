package ilberbot

import "log"

type command func(args ...string) string

var commandMap = map[string]command{}

func RegisterCommand(name string, command command) {
	if _, ok := commandMap[name]; ok {
		log.Fatalf("panic: command '%s' is already registered", name)
	}

	commandMap[name] = command
}

func Dispatch(command string, args ...string) string {
	cmd, ok := commandMap[command]
	if !ok {
		log.Printf("command '%s' not found", command)
		return ""
	}

	return cmd(args...)
}
