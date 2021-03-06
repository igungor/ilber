package command

import (
	"context"

	"github.com/igungor/ilber/bot"
	"github.com/igungor/telegram"
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

func runRayRay(ctx context.Context, b *bot.Bot, msg *telegram.Message) {
	_, err := b.SendMessage(msg.Chat.ID, "malifalitiko!")
	if err != nil {
		b.Logger.Printf("Error while sending message. Err: %v\n", err)
		return
	}
}
