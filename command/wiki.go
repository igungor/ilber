package command

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/igungor/tlbot"
	"golang.org/x/net/context"
)

func init() {
	register(cmdWiki)
}

var cmdWiki = &Command{
	Name:      "bkz",
	ShortLine: "bakınız çok ilginçtir",
	Run:       runWiki,
}

// the best search engine is still google.
// Wikipedia API lacks multi-lingual search.
const wikiURL = "https://ajax.googleapis.com/ajax/services/search/web"

func runWiki(ctx context.Context, b *tlbot.Bot, msg *tlbot.Message) {
	args := msg.Args()
	var txt string
	if len(args) == 0 {
		err := b.SendMessage(msg.Chat.ID, "neye referans vereyim? mesela bana bakın: */bkz İlber Ortaylı*", tlbot.ModeMarkdown, false, nil)
		if err != nil {
			log.Printf("Error while sending message '%v'. Err: %v\n", txt, err)
		}
		return
	}

	qs := strings.Join(args, "+")

	u, _ := url.Parse(wikiURL)
	params := u.Query()
	params.Set("v", "1.0")
	params.Set("q", "wikipedia+"+qs)
	u.RawQuery = params.Encode()

	resp, err := http.Get(u.String())
	if err != nil {
		log.Printf("Error while fetching reference with given criteria '%v'. Err: %v", qs, err)
		return
	}
	defer resp.Body.Close()

	var wikiResult struct {
		ResponseData struct {
			Results []struct {
				URL   string `json:"url"`
				Title string `json:"titleNoFormatting"`
			}
		}
	}

	if err := json.NewDecoder(resp.Body).Decode(&wikiResult); err != nil {
		log.Printf("Error while decoding wiki response: %v\n", err)
		return
	}

	for _, article := range wikiResult.ResponseData.Results {
		if strings.Contains(article.URL, "wikipedia.org/wiki/") {
			articleURL, err := url.QueryUnescape(article.URL)
			if err != nil {
				articleURL = article.URL
			}
			err = b.SendMessage(msg.Chat.ID, articleURL, tlbot.ModeNone, true, nil)
			if err != nil {
				log.Printf("Error while sending message. Err: %v\n", err)
				return
			}
			return
		}
	}

	err = b.SendMessage(msg.Chat.ID, "aradığın referansı bulamadım 🙈", tlbot.ModeMarkdown, true, nil)
	if err != nil {
		log.Printf("Error while sending message. Err: %v\n", err)
		return
	}
}
