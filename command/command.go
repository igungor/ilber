package command

import (
	"strings"
	"sync"

	"github.com/igungor/tlbot"
	"golang.org/x/net/context"
)

// A Command is an implementation of a bot command.
type Command struct {
	// Name of the command without the leading slash.
	Name string

	// Short description of the command
	ShortLine string

	// Hidden enables commands to be hidden from the /help output, such as
	// Telegram's built-in commands and easter eggs.
	Hidden bool

	// Run runs the command.
	Run func(context.Context, *tlbot.Bot, *tlbot.Message)
}

var (
	// mu guards commands-map access
	mu       sync.Mutex
	commands = make(map[string]*Command)
)

func register(cmd *Command) {
	mu.Lock()
	defer mu.Unlock()

	if cmd.Name == "" {
		panic("cannot register command with an empty name")
	}
	if cmd.Run == nil {
		panic("cannot register command with an empty Run function value")
	}
	if _, ok := commands[cmd.Name]; ok {
		panic("plugin already registered: " + cmd.Name)
	}

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
