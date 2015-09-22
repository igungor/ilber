package command

import "github.com/igungor/tlbot"

func init() {
	register(cmdRayRay)
}

var cmdRayRay = &Command{
	Name:      "ray",
	ShortLine: "malifalitiko!",
	Private:   true,
	Run:       runRayRay,
}

func runRayRay(b *tlbot.Bot, msg *tlbot.Message) {
	b.SendMessage(msg.Chat, "malifalitiko!", tlbot.ModeNone, false, nil)
}
