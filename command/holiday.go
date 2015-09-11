package command

import (
	"fmt"
	"sort"
	"time"

	"github.com/igungor/tlbot"
)

func init() {
	sort.Sort(byDate(holidays))
	register(cmdHoliday)
}

var cmdHoliday = &Command{
	Name:      "tatil",
	ShortLine: "ne zaman",
	Run:       runHoliday,
}

var day = 24 * time.Hour

var holidays = []h{
	// 2015
	{"Yılbaşı tatili", newdate("1 Jan 2015"), day},
	{"Çocuk Bayramı", newdate("23 Apr 2015"), day},
	{"İşçi Bayramı", newdate("1 May 2015"), day},
	{"Gençlik Bayramı", newdate("19 May 2015"), day},
	{"Ramazan Bayramı", newdate("18 Jul 2015"), 3 * day},
	{"Zafer Bayramı", newdate("30 Aug 2015"), day},
	{"Kurban Bayramı", newdate("25 Sep 2015"), 4 * day},
	{"Cumhuriyet Bayramı", newdate("29 Oct 2015"), day},
}

type h struct {
	name     string
	date     time.Time
	duration time.Duration
}

type byDate []h

func (d byDate) Len() int           { return len(d) }
func (d byDate) Swap(i, j int)      { d[i], d[j] = d[j], d[i] }
func (d byDate) Less(i, j int) bool { return d[i].date.Before(d[j].date) }

func newdate(date string) time.Time {
	const timeformat = "2 Jan 2006"

	t, _ := time.Parse(timeformat, date)

	return t
}

func in(date, start, end time.Time) bool {
	if date.Equal(start) || date.Equal(end) {
		return true
	}

	return date.After(start) && date.Before(end)
}

func runHoliday(b *tlbot.Bot, msg *tlbot.Message) {
	const timeformat = "2 Jan 2006"
	now := time.Now().UTC()

	for _, t := range holidays {
		if in(now, t.date, t.date.Add(t.duration)) {
			b.SendMessage(msg.Chat, fmt.Sprintf("Bugün %v", t.name), tlbot.ModeNone, false, nil)
			return
		}

		if now.Before(t.date) {
			txt := fmt.Sprintf("En yakın tatil *%v* - %v (*%v* gün)", t.date.Format("_2/01/2006"), t.name, t.duration.Hours()/24)
			b.SendMessage(msg.Chat, txt, tlbot.ModeNone, false, nil)
			return
		}
	}

	b.SendMessage(msg.Chat, "yakın zamanda tatil görünmüyör :(", tlbot.ModeNone, false, nil)
}
