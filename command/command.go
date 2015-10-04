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

const autocorrect = true

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
	if ok {
		return cmd
	}

	if !autocorrect {
		return nil
	}

	//
	// we don't have an exact match. try to guess the input.
	//

	// don't even bother on single letter command inputs
	if len(cmdname) < 2 {
		return nil
	}
	// autocorrect based on levenshtein distance, if possible
	for k := range commands {
		if distance(cmdname, k) <= 2 {
			return commands[k]
		}
	}

	// at least we tried.
	return nil
}

// distance returns the levenshtein distance between given strings.
func distance(s1, s2 string) int {
	var cost, lastdiag, olddiag int
	s1len := len([]rune(s1))
	s2len := len([]rune(s2))

	col := make([]int, s1len+1)

	for y := 1; y <= s1len; y++ {
		col[y] = y
	}

	for x := 1; x <= s2len; x++ {
		col[0] = x
		lastdiag = x - 1
		for y := 1; y <= s1len; y++ {
			olddiag = col[y]
			cost = 0
			if s1[y-1] != s2[x-1] {
				cost = 1
			}
			col[y] = min(col[y]+1, min(col[y-1]+1, lastdiag+cost))
			lastdiag = olddiag
		}
	}
	return col[s1len]
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
