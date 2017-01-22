package command

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/igungor/ilber/bot"
	"github.com/igungor/telegram"
)

func init() {
	register(cmdCandle)
}

var cmdCandle = &Command{
	Name:      "bugunkandilmi",
	ShortLine: "is it candle?",
	Run:       runCandle,
}

var dasCandles = map[string]string{
	// 2016
	"11 Dec 2016": "Mevlid Kandili",

	// 2017
	"30 Mar 2017": "Regaib Kandili",
	"23 Apr 2017": "Mirac Kandili",
	"10 May 2017": "Berat Kandili",
	"21 Jun 2017": "Kadir Gecesi",
	"29 Nov 2017": "Mevlid Kandili",
}

func runCandle(ctx context.Context, b *bot.Bot, msg *telegram.Message) {
	const timeformat = "2 Jan 2006"
	loc, _ := time.LoadLocation("Europe/Istanbul")
	now := time.Now().In(loc).Format(timeformat)
	var txt string
	v, ok := dasCandles[now]
	if !ok {
		txt = "hayır"
	} else {
		txt = fmt.Sprintf("Evet, bugün *%v*\n", v)
	}

	opts := &telegram.SendOptions{ParseMode: telegram.ModeMarkdown}
	_, err := b.SendMessage(msg.Chat.ID, txt, opts)
	if err != nil {
		log.Printf("Error while sending message: %v\n", err)
		return
	}
}
