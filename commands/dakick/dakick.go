package dakick

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/igungor/ilberbot"
)

var (
	near            = "Kadıköy/İstanbul"
	baseURL         = "http://www.google.com/movies?near=" + near
	client          = http.Client{Timeout: 10 * time.Second}
	chromeUserAgent = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2227.0 Safari/537.36"
)

func init() {
	ilberbot.RegisterCommand("/vizyon", movies)
}

func movies(args ...string) string {
	req, _ := http.NewRequest("GET", baseURL, nil)
	req.Header.Set("User-Agent", chromeUserAgent)

	r, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return ""
	}
	defer r.Body.Close()

	doc, err := goquery.NewDocumentFromResponse(r)
	if err != nil {
		log.Println(err)
		return ""
	}

	var buf bytes.Buffer
	buf.WriteString("🎦 Istanbul'da vizyon filmleri\n")

	doc.Find(".theater .desc .name a").Each(func(_ int, s *goquery.Selection) {
		fmt.Println(s.Text())
		if s.Text() == "Cinemaximum Nautilus" {
			s.Closest(".theater").Find(".showtimes .name").Each(func(_ int, sel *goquery.Selection) {
				buf.WriteString(fmt.Sprintf("- %v\n", sel.Text()))
			})
		}
	})

	return buf.String()
}
