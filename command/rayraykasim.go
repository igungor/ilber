package command

import (
	"log"

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

func runRayRay(b *tlbot.Bot, msg *tlbot.Message) {
	err := b.SendMessage(msg.Chat.ID, "malifalitiko!", tlbot.ModeNone, false, nil)
	if err != nil {
		log.Printf("Error while sending message. Err: %v\n", err)
		return
	}
}
