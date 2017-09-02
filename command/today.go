package command

import (
	"context"
	"fmt"
	"time"

	"github.com/igungor/ilber/bot"
	"github.com/igungor/telegram"
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

func runToday(ctx context.Context, b *bot.Bot, msg *telegram.Message) {
	txt := fmt.Sprintf("bugün %v", weekday(time.Now().Weekday()).String())
	_, err := b.SendMessage(msg.Chat.ID, txt)
	if err != nil {
		b.Logger.Printf("Error while sending message. Err: %v\n", err)
	}
}
