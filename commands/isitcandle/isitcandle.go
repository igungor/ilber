package isitcandle

import (
	"fmt"
	"time"

	"github.com/igungor/ilberbot"
)

func init() {
	ilberbot.RegisterCommand("/bugunkandilmi", isitcandle)
	ilberbot.RegisterCommand("/nezamankandil", wheniscandle)
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

// i know it's a lame name but funny nonetheless
func isitcandle(args ...string) string {
	now := time.Now().UTC().Format(timeformat)

	v, ok := dasCandles[now]
	if !ok {
		return "bugun kandil degilmis"
	}

	return fmt.Sprintf("Evet, bugun %v\n", v)

}

func wheniscandle(args ...string) string {
	now := time.Now().UTC()

	for k, v := range dasCandles {
		t, _ := time.Parse(timeformat, k)

		if now.Format(timeformat) == k {
			return fmt.Sprintf("Bugun %v", v)
		}

		if now.Before(t) {
			return fmt.Sprintf("En yakin kandil %v (%v)", v, k)
		}
	}
	return "yakinlarda kandil yok hocam"
}
