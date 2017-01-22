package command

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/igungor/ilber/bot"
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

func runEcho(ctx context.Context, b *bot.Bot, msg *tlbot.Message) {
	args := msg.Args()
	if len(args) == 0 {
		args = []string{"çok cahilsin"}
	}
	txt := fmt.Sprintf("*%v*", strings.Join(args, " "))
	_, err := b.SendMessage(msg.Chat.ID, txt, nil)
	if err != nil {
		log.Printf("Error while sending message: %v\n", err)
		return
	}
}
