package command

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/igungor/tlbot"
)

func init() {
	register(cmdMovies)
}

var cmdMovies = &Command{
	Name:      "vizyon",
	ShortLine: "sinema filan",
	Run:       runMovies,
}

var (
	near            = "Kadıköy/İstanbul"
	movieURL        = "http://www.google.com/movies?near=" + near
	client          = http.Client{Timeout: 10 * time.Second}
	chromeUserAgent = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2227.0 Safari/537.36"
)

func runMovies(b *tlbot.Bot, msg *tlbot.Message) {
	req, _ := http.NewRequest("GET", movieURL, nil)
	req.Header.Set("User-Agent", chromeUserAgent)

	r, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	defer r.Body.Close()

	doc, err := goquery.NewDocumentFromResponse(r)
	if err != nil {
		log.Printf("(movie) Error while fetching DOM: %v\n", err)
		return
	}

	var buf bytes.Buffer
	buf.WriteString("🎦 Istanbul'da vizyon filmleri\n")

	doc.Find(".theater .desc .name a").Each(func(_ int, s *goquery.Selection) {
		if s.Text() == "Cinemaximum Nautilus" {
			s.Closest(".theater").Find(".showtimes .name").Each(func(_ int, sel *goquery.Selection) {
				buf.WriteString(fmt.Sprintf("- %v\n", sel.Text()))
			})
		}
	})

	b.SendMessage(msg.From, buf.String(), tlbot.ModeNone, false, nil)
}
