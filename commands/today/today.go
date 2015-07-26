package today

import (
	"fmt"
	"time"

	"github.com/igungor/ilberbot"
)

func init() {
	ilberbot.RegisterCommand("/bugun", today)
}

type weekday time.Weekday

var days = [...]string{
	"pazar",
	"pazartesi",
	"sali",
	"carsamba",
	"persembe",
	"cuma",
	"cumartesi",
}

func (w weekday) String() string {
	return days[w]
}

func today(args ...string) string {
	return fmt.Sprintf("bugun %v", weekday(time.Now().Weekday()).String())
}
