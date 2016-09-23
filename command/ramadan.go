package command

import (
	"fmt"
	"log"
	"time"

	"github.com/igungor/tlbot"
	"golang.org/x/net/context"
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
	"18 Jun 2015": {"20:48", "03:22"},
	"19 Jun 2015": {"20:48", "03:22"},
	"20 Jun 2015": {"20:48", "03:22"},
	"21 Jun 2015": {"20:48", "03:22"},
	"22 Jun 2015": {"20:49", "03:23"},
	"23 Jun 2015": {"20:49", "03:23"},
	"24 Jun 2015": {"20:49", "03:23"},
	"25 Jun 2015": {"20:49", "03:23"},
	"26 Jun 2015": {"20:49", "03:24"},
	"27 Jun 2015": {"20:49", "03:24"},
	"28 Jun 2015": {"20:49", "03:25"},
	"29 Jun 2015": {"20:49", "03:26"},
	"30 Jun 2015": {"20:49", "03:26"},
	"1 Jul 2015":  {"20:49", "03:27"},
	"2 Jul 2015":  {"20:49", "03:28"},
	"3 Jul 2015":  {"20:49", "03:28"},
	"4 Jul 2015":  {"20:49", "03:29"},
	"5 Jul 2015":  {"20:48", "03:30"},
	"6 Jul 2015":  {"20:48", "03:31"},
	"7 Jul 2015":  {"20:48", "03:32"},
	"8 Jul 2015":  {"20:48", "03:33"},
	"9 Jul 2015":  {"20:47", "03:35"},
	"10 Jul 2015": {"20:47", "03:35"},
	"11 Jul 2015": {"20:46", "03:37"},
	"12 Jul 2015": {"20:46", "03:38"},
	"13 Jul 2015": {"20:46", "03:39"},
	"14 Jul 2015": {"20:45", "03:40"},
	"15 Jul 2015": {"20:44", "03:41"},
	"16 Jul 2015": {"20:44", "03:43"},

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
}

func populate() {
	loc, _ := time.LoadLocation(timezone)

	for k, v := range _callTime {
		iftar, _ := time.ParseInLocation(timeFormatLong, fmt.Sprintf("%v %v", k, v.iftar), loc)
		sahur, _ := time.ParseInLocation(timeFormatLong, fmt.Sprintf("%v %v", k, v.sahur), loc)

		callTime[k] = timePair{iftar, sahur}
	}
}

func runPrayerCall(ctx context.Context, b *tlbot.Bot, msg *tlbot.Message) {
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

	err := b.SendMessage(msg.Chat.ID, txt, tlbot.ModeNone, false, nil)
	if err != nil {
		log.Printf("Error while sending message. Err: %v\n", err)
		return
	}
}

func runFoodFast(ctx context.Context, b *tlbot.Bot, msg *tlbot.Message) {
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

	err := b.SendMessage(msg.Chat.ID, txt, tlbot.ModeNone, false, nil)
	if err != nil {
		log.Printf("Error while sending message. Err: %v\n", err)
		return
	}
}

func runFoodDawn(ctx context.Context, b *tlbot.Bot, msg *tlbot.Message) {
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

	err := b.SendMessage(msg.Chat.ID, txt, tlbot.ModeNone, false, nil)
	if err != nil {
		log.Printf("Error while sending message. Err: %v\n", err)
		return
	}
}
