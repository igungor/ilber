package command

import (
	"sync"

	"github.com/igungor/tlbot"
)

var mu sync.Mutex

// A Command is an implementation of a bot command.
type Command struct {
	// Name of the command
	Name string

	// Run runs the command.
	// The args are the arguments after the command name.
	Run func(bot *tlbot.Bot, msg *tlbot.Message)
}

var commands = make(map[string]*Command)

func register(cmd *Command) {
	mu.Lock()
	defer mu.Unlock()

	commands[cmd.Name] = cmd
}

// Lookup looks-up name from registered commands store and returns
// corresponding Command, if any.
func Lookup(name string) *Command {
	mu.Lock()
	defer mu.Unlock()

	cmd, ok := commands[name]
	if !ok {
		return nil
	}
	return cmd
}
