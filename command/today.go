package command

import (
	"context"
	"fmt"
	"log"
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
	"salı",
	"çarşamba",
	"perşembe",
	"cuma",
	"cumartesi",
}

func (w weekday) String() string {
	return days[w]
}

func runToday(ctx context.Context, b *tlbot.Bot, msg *tlbot.Message) {
	txt := fmt.Sprintf("bugün %v", weekday(time.Now().Weekday()).String())
	_, err := b.SendMessage(msg.Chat.ID, txt, nil)
	if err != nil {
		log.Printf("Error while sending message. Err: %v\n", err)
	}
}
