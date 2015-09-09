package command

import (
	"fmt"
	"strings"

	"github.com/igungor/tlbot"
)

func init() {
	register(cmdEcho)
}

var cmdEcho = &Command{
	Name:      "echo",
	ShortLine: "çok cahilsin",
	Run:       runEcho,
}

func runEcho(b *tlbot.Bot, msg *tlbot.Message) {
	args := msg.Args()
	txt := fmt.Sprintf("*%v*", strings.Join(args, " "))
	b.SendMessage(msg.From, txt, tlbot.ModeMarkdown, false, nil)
}
