package command

import (
	"log"
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
	if len(args) == 0 {
		err := b.SendMessage(msg.Chat.ID, "neye referans vereyim? mesela bana bakın: */bkz İlber Ortaylı*", tlbot.ModeMarkdown, false, nil)
		if err != nil {
			log.Printf("Error while sending message. Err: %v\n", err)
		}
		return
	}

	googleAPIKey := ctx.Value("googleAPIKey").(string)
	searchEngineID := ctx.Value("googleSearchEngineID").(string)

	terms := []string{"wikipedia"}
	terms = append(terms, args...)

	urls, err := search(googleAPIKey, searchEngineID, "", terms...)
	if err != nil {
		log.Printf("Error while 'bkz' query. Err: %v\n", err)
		if err == errSearchQuotaExceeded {
			b.SendMessage(msg.Chat.ID, `¯\_(ツ)_/¯`, tlbot.ModeNone, false, nil)
		}
		return
	}

	for _, articleURL := range urls {
		if strings.Contains(articleURL, "wikipedia.org/wiki/") {
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
