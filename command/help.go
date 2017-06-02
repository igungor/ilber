package command

import (
	"bytes"
	"context"
	"fmt"
	"sort"

	"github.com/igungor/ilber/bot"
	"github.com/igungor/telegram"
)

func init() {
	register(cmdHelp)
}

var cmdHelp = &Command{
	Name:      "help",
	ShortLine: "imdaat!",
	Run:       runHelp,
}

func runHelp(ctx context.Context, b *bot.Bot, msg *telegram.Message) {
	b.SendMessage(msg.Chat.ID, help())
}

type byName []*Command

func (b byName) Len() int           { return len(b) }
func (b byName) Less(i, j int) bool { return b[i].Name < b[j].Name }
func (b byName) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }

func help() string {
	var cmds []*Command
	for _, cmd := range commands {
		cmds = append(cmds, cmd)
	}

	sort.Sort(byName(cmds))

	var buf bytes.Buffer
	buf.WriteString("şunlar var:\n\n")
	for _, cmd := range cmds {
		// do not include hidden commands
		if cmd.Hidden {
			continue
		}
		buf.WriteString(fmt.Sprintf("/%v - %v\n", cmd.Name, cmd.ShortLine))
	}

	return buf.String()
}
