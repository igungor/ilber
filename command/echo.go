package command

import (
	"fmt"
	"log"
	"strings"

	"github.com/igungor/tlbot"
	"golang.org/x/net/context"
)

func init() {
	register(cmdEcho)
}

var cmdEcho = &Command{
	Name:      "echo",
	ShortLine: "çok cahilsin",
	Run:       runEcho,
}

func runEcho(ctx context.Context, b *tlbot.Bot, msg *tlbot.Message) {
	args := msg.Args()
	if len(args) == 0 {
		args = []string{"çok cahilsin"}
	}
	txt := fmt.Sprintf("*%v*", strings.Join(args, " "))
	err := b.SendMessage(msg.Chat.ID, txt, tlbot.ModeMarkdown, false, nil)
	if err != nil {
		log.Printf("Error while sending message: %v\n", err)
		return
	}
}
