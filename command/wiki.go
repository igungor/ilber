package command

import (
	"encoding/xml"
	"log"
	"net/url"
	"strings"

	"github.com/igungor/tlbot"
)

func init() {
	register(cmdWiki)
}

const wikiURL = "https://en.wikipedia.org/w/api.php"

var cmdWiki = &Command{
	Name:      "bkz",
	ShortLine: "bakiniz cok ilginctir",
	Run:       runWiki,
}

func runWiki(b *tlbot.Bot, msg *tlbot.Message) {
	args := msg.Args()
	var txt string
	if len(args) == 0 {
		err := b.SendMessage(msg.Chat, "neye referans vereyim? mesela bana bakın: */bkz Ilber Ortayli*", tlbot.ModeMarkdown, false, nil)
		if err != nil {
			log.Printf("[wiki] Error while sending message '%v'. Err: %v\n", txt, err)
		}
		return
	}

	u, err := url.Parse(wikiURL)
	if err != nil {
		log.Printf("[wiki] Error while parsing URL '%v'. Err: %v\n", wikiURL, err)
		return
	}

	qs := strings.Join(args, " ")
	params := u.Query()
	params.Set("action", "opensearch")
	params.Set("format", "xml")
	params.Set("search", qs)
	u.RawQuery = params.Encode()

	resp, err := httpclient.Get(u.String())
	if err != nil {
		log.Printf("[wiki] Error while fetching article for query '%v'. Err: %v\n", qs, err)
		return
	}
	defer resp.Body.Close()

	var res result
	err = xml.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		log.Printf("[wiki] Error while decoding response '%v'. Err: %v\n", res, err)
		return
	}

	if len(res.Section.Items) == 0 {
		b.SendMessage(msg.Chat, "aradığın referansı bulamadım 🙈", tlbot.ModeMarkdown, true, nil)
		return
	}

	article := res.Section.Items[0]
	txt = article.URL
	err = b.SendMessage(msg.Chat, txt, tlbot.ModeNone, true, nil)
	if err != nil {
		log.Printf("[wiki] Error while sending message to Telegram servers. Err: '%v'\nMessage: %v", err, txt)
		return
	}
}

// wikipedia api response
type result struct {
	XMLName xml.Name `xml:"SearchSuggestion"`
	Section struct {
		Items []struct {
			Title       string `xml:"Text"`
			URL         string `xml:"Url"`
			Description string
			Image       string
		} `xml:"Item"`
	} `xml:"Section"`
}
