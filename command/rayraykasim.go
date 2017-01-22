package command

import (
	"context"
	"log"

	"github.com/igungor/ilber/bot"
	"github.com/igungor/tlbot"
)

func init() {
	register(cmdRayRay)
}

var cmdRayRay = &Command{
	Name:      "ray",
	ShortLine: "malifalitiko!",
	Hidden:    true,
	Run:       runRayRay,
}

func runRayRay(ctx context.Context, b *bot.Bot, msg *tlbot.Message) {
	_, err := b.SendMessage(msg.Chat.ID, "malifalitiko!", nil)
	if err != nil {
		log.Printf("Error while sending message. Err: %v\n", err)
		return
	}
}
