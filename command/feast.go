package command

import (
	"fmt"
	"time"

	"github.com/igungor/tlbot"
)

func init() {
	// XXX: enable them when the time comes.
	// populate()
	// register(cmdPrayerCall)
	// register(cmdFoodFast)
	// register(cmdFoodDawn)
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
}

func populate() {
	loc, _ := time.LoadLocation("Europe/Istanbul")

	for k, v := range _callTime {
		iftar, _ := time.ParseInLocation(timeFormatLong, fmt.Sprintf("%v %v", k, v.iftar), loc)
		sahur, _ := time.ParseInLocation(timeFormatLong, fmt.Sprintf("%v %v", k, v.sahur), loc)

		callTime[k] = timePair{iftar, sahur}
	}
}

func runPrayerCall(b *tlbot.Bot, msg *tlbot.Message) {
	loc, _ := time.LoadLocation(timezone)

	now := time.Now().In(loc)
	nowstr := now.Format(timeFormat)

	timepair, ok := callTime[nowstr]
	if !ok {
		b.SendMessage(msg.Chat, "galiba oruç bitti", tlbot.ModeNone, false, nil)
		return
	}

	if now.After(timepair.iftar) {
		b.SendMessage(msg.Chat, "okundu", tlbot.ModeNone, false, nil)
		return
	}

	if now.Before(timepair.sahur) {
		b.SendMessage(msg.Chat, "sahur henüz okunmadı", tlbot.ModeNone, false, nil)
		return
	}

	if now.After(timepair.sahur) && now.Hour() < 6 {
		b.SendMessage(msg.Chat, "sahur okundu ama iftara daha çok var", tlbot.ModeNone, false, nil)
		return
	}

	// after sahur and before iftar, hence NO
	b.SendMessage(msg.Chat, randChoice(noes), tlbot.ModeNone, false, nil)
}

func runFoodFast(b *tlbot.Bot, msg *tlbot.Message) {
	loc, _ := time.LoadLocation(timezone)

	now := time.Now().In(loc)
	nowstr := now.Format(timeFormat)

	timepair, ok := callTime[nowstr]
	if !ok {
		b.SendMessage(msg.Chat, "galiba oruç bitti", tlbot.ModeNone, false, nil)
		return
	}

	b.SendMessage(msg.Chat, timepair.iftar.Format("15:04"), tlbot.ModeNone, false, nil)
}

func runFoodDawn(b *tlbot.Bot, msg *tlbot.Message) {
	loc, _ := time.LoadLocation(timezone)

	now := time.Now().In(loc)
	nowstr := now.Format(timeFormat)

	timepair, ok := callTime[nowstr]
	if !ok {
		b.SendMessage(msg.Chat, "galiba oruç bitti", tlbot.ModeNone, false, nil)
		return
	}

	b.SendMessage(msg.Chat, timepair.iftar.Format("15:04"), tlbot.ModeNone, false, nil)
}
