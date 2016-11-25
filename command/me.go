package command

import (
	"fmt"
	"log"
	"strings"

	"github.com/igungor/tlbot"
	"golang.org/x/net/context"
)

func init() {
	register(cmdMe)
}

var cmdMe = &Command{
	Name:      "me",
	ShortLine: "ay resmen ay-ar-sii",
	Run:       runMe,
}

func runMe(ctx context.Context, b *tlbot.Bot, msg *tlbot.Message) {
	args := msg.Args()
	if len(args) == 0 {
		args = []string{"hmmmmm"}
	}

	user := msg.From.FirstName
	if user == "" {
		user = msg.From.Username
	}

	txt := fmt.Sprintf("`* %v %v`", user, strings.Join(args, " "))
	err := b.SendMessage(msg.Chat.ID, txt, tlbot.ModeMarkdown, false, nil)
	if err != nil {
		log.Printf("Error while sending message: %v\n", err)
		return
	}
}
