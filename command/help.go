package command

import (
	"bytes"
	"fmt"
	"sort"

	"github.com/igungor/tlbot"
)

func init() {
	register(cmdHelp)
}

var cmdHelp = &Command{
	Name:      "help",
	ShortLine: "imdaat!",
	Run:       runHelp,
}

func runHelp(b *tlbot.Bot, msg *tlbot.Message) {
	b.SendMessage(msg.Chat, help(), tlbot.ModeNone, false, nil)
}

type byName []*Command

func (b byName) Len() int           { return len(b) }
func (b byName) Less(i, j int) bool { return b[i].Name < b[j].Name }
func (b byName) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }

func help() string {
	var buf bytes.Buffer

	var cmds []*Command
	for _, cmd := range commands {
		cmds = append(cmds, cmd)
	}

	sort.Sort(byName(cmds))

	buf.WriteString("şunlar var:\n\n")
	for _, cmd := range cmds {
		// do not include private commands
		if cmd.Private {
			continue
		}
		buf.WriteString(fmt.Sprintf("/%v - %v\n", cmd.Name, cmd.ShortLine))
	}

	return buf.String()
}
