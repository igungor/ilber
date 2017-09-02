package command

import (
	"context"
	"fmt"
	"time"

	"github.com/igungor/ilber/bot"
	"github.com/igungor/telegram"
)

func init() {
	populate()
	register(cmdPrayerCall)
	register(cmdFoodFast)
	register(cmdFoodDawn)
}

var (
	cmdPrayerCall = &Command{
		Name:      "okundumu",
		ShortLine: "is it read?",
		Run:       runPrayerCall,
	}
	cmdFoodFast = &Command{
		Name:      "iftar",
		ShortLine: "seninle iftar ediyorum",
		Run:       runFoodFast,
	}
	cmdFoodDawn = &Command{
		Name:      "sahur",
		ShortLine: "sahura kalkti mi?",
		Run:       runFoodDawn,
	}
)

const (
	timeFormat     = "2 Jan 2006"
	timeFormatLong = "2 Jan 2006 15:04"
	timezone       = "Europe/Istanbul"
)

var noes = []string{
	"hayır",
	"hayır.",
	"Hayır",
	"Hayır.",
	"no",
	"no.",
	"NO.",
	"NO!",
	"NO",
	"nope",
	"nicht",
	"okunmadı",
	"hayır okunmadı",
}

type timePair struct {
	iftar time.Time
	sahur time.Time
}

var callTime = map[string]timePair{}

var _callTime = map[string]struct {
	iftar, sahur string
}{
	"6 Jun 2016":  {"20:41", "03:26"},
	"7 Jun 2016":  {"20:42", "03:26"},
	"8 Jun 2016":  {"20:43", "03:25"},
	"9 Jun 2016":  {"20:43", "03:25"},
	"10 Jun 2016": {"20:44", "03:24"},
	"11 Jun 2016": {"20:44", "03:24"},
	"12 Jun 2016": {"20:45", "03:23"},
	"13 Jun 2016": {"20:45", "03:23"},
	"14 Jun 2016": {"20:46", "03:23"},
	"15 Jun 2016": {"20:46", "03:23"},
	"16 Jun 2016": {"20:46", "03:23"},
	"17 Jun 2016": {"20:47", "03:23"},
	"18 Jun 2016": {"20:47", "03:23"},
	"19 Jun 2016": {"20:47", "03:23"},
	"20 Jun 2016": {"20:48", "03:23"},
	"21 Jun 2016": {"20:48", "03:23"},
	"22 Jun 2016": {"20:48", "03:23"},
	"23 Jun 2016": {"20:48", "03:24"},
	"24 Jun 2016": {"20:48", "03:24"},
	"25 Jun 2016": {"20:49", "03:24"},
	"26 Jun 2016": {"20:49", "03:25"},
	"27 Jun 2016": {"20:49", "03:25"},
	"28 Jun 2016": {"20:49", "03:26"},
	"29 Jun 2016": {"20:49", "03:27"},
	"30 Jun 2016": {"20:49", "03:27"},
	"1 Jul 2016":  {"20:48", "03:28"},
	"2 Jul 2016":  {"20:48", "03:29"},
	"3 Jul 2016":  {"20:48", "03:30"},
	"4 Jul 2016":  {"20:48", "03:31"},

	"27 May 2017": {"20:33", "03:36"},
	"28 May 2017": {"20:34", "03:35"},
	"29 May 2017": {"20:35", "03:34"},
	"30 May 2017": {"20:36", "03:33"},
	"31 May 2017": {"20:37", "03:32"},
	"1 Jun 2017":  {"20:37", "03:32"},
	"2 Jun 2017":  {"20:38", "03:30"},
	"3 Jun 2017":  {"20:39", "03:29"},
	"4 Jun 2017":  {"20:40", "03:28"},
	"5 Jun 2017":  {"20:40", "03:28"},
	"6 Jun 2017":  {"20:41", "03:27"},
	"7 Jun 2017":  {"20:42", "03:26"},
	"8 Jun 2017":  {"20:42", "03:26"},
	"9 Jun 2017":  {"20:43", "03:25"},
	"10 Jun 2017": {"20:43", "03:25"},
	"11 Jun 2017": {"20:44", "03:24"},
	"12 Jun 2017": {"20:44", "03:24"},
	"13 Jun 2017": {"20:45", "03:24"},
	"14 Jun 2017": {"20:45", "03:23"},
	"15 Jun 2017": {"20:46", "03:23"},
	"16 Jun 2017": {"20:46", "03:23"},
	"17 Jun 2017": {"20:46", "03:23"},
	"18 Jun 2017": {"20:47", "03:23"},
	"19 Jun 2017": {"20:47", "03:23"},
	"20 Jun 2017": {"20:47", "03:23"},
	"21 Jun 2017": {"20:48", "03:23"},
	"22 Jun 2017": {"20:48", "03:24"},
	"23 Jun 2017": {"20:48", "03:24"},
	"24 Jun 2017": {"20:48", "03:24"},
}

func populate() {
	loc, _ := time.LoadLocation(timezone)

	for k, v := range _callTime {
		iftar, _ := time.ParseInLocation(timeFormatLong, fmt.Sprintf("%v %v", k, v.iftar), loc)
		sahur, _ := time.ParseInLocation(timeFormatLong, fmt.Sprintf("%v %v", k, v.sahur), loc)

		callTime[k] = timePair{iftar, sahur}
	}
}

func runPrayerCall(ctx context.Context, b *bot.Bot, msg *telegram.Message) {
	loc, _ := time.LoadLocation(timezone)

	now := time.Now().In(loc)
	nowstr := now.Format(timeFormat)

	timepair, ok := callTime[nowstr]
	var txt string
	switch {
	case !ok:
		txt = "galiba oruç bitti"
	case now.After(timepair.iftar):
		txt = "okundu"
	case now.After(timepair.sahur) && now.Hour() < 6:
		txt = "sahur okundu ama iftara daha çok var"
	case now.Before(timepair.sahur):
		txt = "sahur henüz okunmadı"
	default:
		// after sahur and before iftar, hence NO
		txt = randChoice(noes)
	}

	_, err := b.SendMessage(msg.Chat.ID, txt)
	if err != nil {
		b.Logger.Printf("Error while sending message. Err: %v\n", err)
		return
	}
}

func runFoodFast(ctx context.Context, b *bot.Bot, msg *telegram.Message) {
	loc, _ := time.LoadLocation(timezone)

	now := time.Now().In(loc)
	nowstr := now.Format(timeFormat)

	timepair, ok := callTime[nowstr]
	var txt string
	if ok {
		txt = timepair.iftar.Format("15:04")
	} else {
		txt = "galiba oruç bitti"
	}

	_, err := b.SendMessage(msg.Chat.ID, txt)
	if err != nil {
		b.Logger.Printf("Error while sending message. Err: %v\n", err)
		return
	}
}

func runFoodDawn(ctx context.Context, b *bot.Bot, msg *telegram.Message) {
	loc, _ := time.LoadLocation(timezone)

	now := time.Now().In(loc)
	nowstr := now.Format(timeFormat)

	timepair, ok := callTime[nowstr]
	var txt string
	if ok {
		txt = timepair.sahur.Format("15:04")
	} else {
		txt = "galiba oruç bitti"
	}

	_, err := b.SendMessage(msg.Chat.ID, txt)
	if err != nil {
		b.Logger.Printf("Error while sending message. Err: %v\n", err)
		return
	}
}
