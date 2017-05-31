package command

import (
	"context"
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/igungor/ilber/bot"
	"github.com/igungor/telegram"
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

var holidays = []holiday{
	// 2017
	{"Yılbaşı tatili", newdate("1 Jan 2017"), day},
	{"Çocuk Bayramı", newdate("23 Apr 2017"), day},
	{"İşçi Bayramı", newdate("1 May 2017"), day},
	{"Gençlik Bayramı", newdate("19 May 2017"), day},
	{"Ramazan Bayramı", newdate("24 Jun 2017"), 3 * day},
	{"Demokrasi Şeysi", newdate("15 Jun 2017"), day},
	{"Zafer Bayramı", newdate("30 Aug 2017"), day},
	{"Kurban Bayramı", newdate("31 Aug 2017"), 12 * time.Hour},
	{"Kurban Bayramı", newdate("1 Sep 2017"), 4 * day},
	{"Cumhuriyet Bayramı", newdate("29 Oct 2017"), day + 12*time.Hour},
}

type holiday struct {
	name     string
	date     time.Time
	duration time.Duration
}

type byDate []holiday

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

func runHoliday(ctx context.Context, b *bot.Bot, msg *telegram.Message) {
	now := time.Now().UTC()
	md := telegram.WithParseMode(telegram.ModeMarkdown)
	for _, holiday := range holidays {
		if in(now, holiday.date, holiday.date.Add(holiday.duration)) {
			_, err := b.SendMessage(msg.Chat.ID, fmt.Sprintf("Bugün %v", holiday.name), md)
			if err != nil {
				log.Printf("Error while sending message. Err: %v\n", err)
			}
			return
		}

		if now.Before(holiday.date) {
			txt := fmt.Sprintf("En yakın tatil *%v* - %v (*%v* gün)", holiday.date.Format("_2/01/2006"), holiday.name, holiday.duration.Hours()/24)
			_, err := b.SendMessage(msg.Chat.ID, txt, md)
			if err != nil {
				log.Printf("Error while sending message. Err: %v\n", err)
			}
			return
		}
	}

	_, err := b.SendMessage(msg.Chat.ID, "yakın zamanda tatil görünmüyor :(", md)
	if err != nil {
		log.Printf("Error while sending message. Err: %v\n", err)
		return
	}
}
