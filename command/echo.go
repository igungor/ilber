package command

import (
	"context"
	"fmt"
	"log"
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

	opts := telegram.SendOptions{
		ParseMode: telegram.ModeMarkdown,
	}
	txt := fmt.Sprintf("*%v*", strings.Join(args, " "))
	_, err := b.SendMessage(msg.Chat.ID, txt, &opts)
	if err != nil {
		log.Printf("Error while sending message: %v\n", err)
		return
	}
}
