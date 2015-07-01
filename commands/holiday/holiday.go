package holiday

import (
	"fmt"
	"sort"
	"time"

	"github.com/igungor/ilberbot"
)

func init() {
	sort.Sort(byDate(holidays))

	ilberbot.RegisterCommand("/tatil", holiday)
}

var day = 24 * time.Hour

var holidays = []h{
	// 2015
	{"Yilbasi tatili", newdate("1 Jan 2015"), day},
	{"Cocuk Bayrami", newdate("23 Apr 2015"), day},
	{"Isci Bayrami", newdate("1 May 2015"), day},
	{"Genclik Bayrami", newdate("19 May 2015"), day},
	{"Ramazan Bayrami", newdate("18 Jul 2015"), 3 * day},
	{"Zafer Bayrami", newdate("30 Aug 2015"), day},
	{"Kurban Bayrami", newdate("25 Sep 2015"), 4 * day},
	{"Cumhuriyet Bayrami", newdate("29 Oct 2015"), day},
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

func holiday(args ...string) string {
	const timeformat = "2 Jan 2006"

	now := time.Now().UTC()

	for _, t := range holidays {
		if in(now, t.date, t.date.Add(t.duration)) {
			return fmt.Sprintf("Bugun %v", t.name)
		}

		if now.Before(t.date) {
			return fmt.Sprintf(
				"En yakin tatil zamani %v - %v (%v gun)",
				t.date.Format("_2 Jan"), t.name, t.duration.String(),
			)
		}
	}

	return "yakinlarda tatil gorunmuyor :("
}
