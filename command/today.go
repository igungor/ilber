package command

import (
	"fmt"
	"time"

	"github.com/igungor/tlbot"
)

func init() {
	register(cmdToday)
}

var cmdToday = &Command{
	Name:      "bugun",
	ShortLine: "bugün günlerden ne?",
	Run:       runToday,
}

type weekday time.Weekday

var days = [...]string{
	"pazar",
	"pazartesi",
	"sali",
	"carsamba",
	"persembe",
	"cuma",
	"cumartesi",
}

func (w weekday) String() string {
	return days[w]
}

func runToday(b *tlbot.Bot, msg *tlbot.Message) {
	txt := fmt.Sprintf("bugün %v", weekday(time.Now().Weekday()).String())
	b.SendMessage(msg.Chat, txt, tlbot.ModeNone, false, nil)
}
