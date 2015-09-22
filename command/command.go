package command

import (
	"strings"
	"sync"

	"github.com/igungor/tlbot"
)

// A Command is an implementation of a bot command.
type Command struct {
	// Name of the command
	Name string

	// Short description of the command
	ShortLine string

	// Some commands would like to stay private, such as easter eggs or
	// built-in commands. Respect to their choices.
	Private bool

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

// Lookup looks-up name from registered commands and returns
// corresponding Command if any.
func Lookup(cmdname string) *Command {
	mu.Lock()
	defer mu.Unlock()

	cmdname = strings.TrimSuffix(cmdname, "@ilberbot")
	cmd, ok := commands[cmdname]
	if !ok {
		return nil
	}
	return cmd
}
