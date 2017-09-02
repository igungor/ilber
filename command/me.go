package command

import (
	"context"
	"fmt"
	"strings"

	"github.com/igungor/ilber/bot"
	"github.com/igungor/telegram"
)

func init() {
	register(cmdMe)
}

var cmdMe = &Command{
	Name:      "me",
	ShortLine: "ay resmen ay-ar-sii",
	Run:       runMe,
}

func runMe(ctx context.Context, b *bot.Bot, msg *telegram.Message) {
	args := msg.Args()
	if len(args) == 0 {
		args = []string{"hmmmmm"}
	}

	user := msg.From.FirstName
	if user == "" {
		user = msg.From.Username
	}

	txt := fmt.Sprintf("`* %v %v`", user, strings.Join(args, " "))
	md := telegram.WithParseMode(telegram.ModeMarkdown)
	_, err := b.SendMessage(msg.Chat.ID, txt, md)
	if err != nil {
		b.Logger.Printf("Error while sending message: %v\n", err)
		return
	}
}
