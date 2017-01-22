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
	register(cmdMe)
}

var cmdMe = &Command{
	Name:      "me",
	ShortLine: "ay resmen ay-ar-sii",
	Run:       runMe,
}

func runMe(ctx context.Context, b *bot.Bot, msg *tlbot.Message) {
	args := msg.Args()
	if len(args) == 0 {
		args = []string{"hmmmmm"}
	}

	user := msg.From.FirstName
	if user == "" {
		user = msg.From.Username
	}

	txt := fmt.Sprintf("`* %v %v`", user, strings.Join(args, " "))
	opts := &tlbot.SendOptions{ParseMode: tlbot.ModeMarkdown}
	_, err := b.SendMessage(msg.Chat.ID, txt, opts)
	if err != nil {
		log.Printf("Error while sending message: %v\n", err)
		return
	}
}
