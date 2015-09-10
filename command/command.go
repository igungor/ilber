package command

import (
	"sync"

	"github.com/igungor/tlbot"
)

// A Command is an implementation of a bot command.
type Command struct {
	// Name of the command
	Name string

	// Short description of the command
	ShortLine string

	// Run runs the command.
	Run func(bot *tlbot.Bot, msg *tlbot.Message)
}

var (
	mu       sync.Mutex
	commands = make(map[string]*Command)
)

func register(cmd *Command) {
	mu.Lock()
	defer mu.Unlock()

	commands[cmd.Name] = cmd
}

// Lookup looks-up name from registered commands store and returns
// corresponding Command if any.
func Lookup(name string) *Command {
	mu.Lock()
	defer mu.Unlock()

	cmd, ok := commands[name]
	if !ok {
		return nil
	}
	return cmd
}
