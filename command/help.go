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

type byName []*Command

func (b byName) Len() int           { return len(b) }
func (b byName) Less(i, j int) bool { return b[i].Name < b[j].Name }
func (b byName) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }

func runHelp(b *tlbot.Bot, msg *tlbot.Message) {
	var buf bytes.Buffer

	var cmds []*Command
	for _, cmd := range commands {
		cmds = append(cmds, cmd)
	}

	sort.Sort(byName(cmds))

	buf.WriteString("şunlar var:\n")
	for _, cmd := range cmds {
		buf.WriteString(fmt.Sprintf("*%v* - %v\n", cmd.Name, cmd.ShortLine))
	}
	b.SendMessage(msg.From, buf.String(), tlbot.ModeMarkdown, false, nil)
}
