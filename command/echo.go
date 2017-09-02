package command

import (
	"context"
	"fmt"
	"strings"

	"github.com/igungor/ilber/bot"
	"github.com/igungor/telegram"
)

func init() {
	register(cmdEcho)
}

var cmdEcho = &Command{
	Name:      "echo",
	ShortLine: "çok cahilsin",
	Run:       runEcho,
}

func runEcho(ctx context.Context, b *bot.Bot, msg *telegram.Message) {
	args := msg.Args()
	if len(args) == 0 {
		args = []string{"çok cahilsin"}
	}

	txt := fmt.Sprintf("*%v*", strings.Join(args, " "))
	_, err := b.SendMessage(msg.Chat.ID, txt, telegram.WithParseMode(telegram.ModeMarkdown))
	if err != nil {
		b.Logger.Printf("Error while sending message: %v\n", err)
		return
	}
}
