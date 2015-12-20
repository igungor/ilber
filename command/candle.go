package command

import (
	"fmt"
	"log"
	"time"

	"github.com/igungor/tlbot"
	"golang.org/x/net/context"
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
	"2 Jan 2015":  "Mevlid Kandili",
	"23 Apr 2015": "Regaib Kandili",
	"15 May 2015": "Mirac Kandili",
	"1 Jun 2015":  "Berat Kandili",
	"13 Jul 2015": "Kadir Gecesi",
	"22 Dec 2015": "Mevlid Kandili",

	"7 Apr 2016":  "Regaib Kandili",
	"3 Apr 2016":  "Mirac Kandili",
	"21 May 2016": "Berat Kandili",
	"1 Jul 2016":  "Kadir Gecesi",
	"11 Dec 2016": "Mevlid Kandili",
}

func runCandle(ctx context.Context, b *tlbot.Bot, msg *tlbot.Message) {
	const timeformat = "2 Jan 2006"
	var txt string
	now := time.Now().UTC().Format(timeformat)
	v, ok := dasCandles[now]
	if !ok {
		txt = "hayır"
	} else {
		txt = fmt.Sprintf("Evet, bugün *%v*\n", v)
	}

	err := b.SendMessage(msg.Chat.ID, txt, tlbot.ModeMarkdown, false, nil)
	if err != nil {
		log.Printf("Error while sending message: %v\n", err)
		return
	}
}
