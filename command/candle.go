package command

import (
	"fmt"
	"time"

	"github.com/igungor/tlbot"
)

func init() {
	register(cmdCandle)
}

var cmdCandle = &Command{
	Name:      "bugunkandilmi",
	ShortLine: "is it candle?",
	Run:       runCandle,
}

const timeformat = "2 Jan 2006"

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

func runCandle(b *tlbot.Bot, msg *tlbot.Message) {
	var txt string
	now := time.Now().UTC().Format(timeformat)
	v, ok := dasCandles[now]
	if !ok {
		txt = "hayır"
	} else {
		txt = fmt.Sprintf("Evet, bugün *%v*\n", v)
	}

	b.SendMessage(msg.From, txt, tlbot.ModeMarkdown, false, nil)
}
