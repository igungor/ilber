package command

import (
	"context"
	"log"
	"strings"

	"github.com/igungor/tlbot"
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
	opts := &tlbot.SendOptions{ParseMode: tlbot.ModeMarkdown}
	if len(args) == 0 {
		txt := "neye referans vereyim? mesela bana bakın: */bkz İlber Ortaylı*"
		_, err := b.SendMessage(msg.Chat.ID, txt, opts)
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
			b.SendMessage(msg.Chat.ID, `¯\_(ツ)_/¯`, nil)
		}
		return
	}

	for _, articleURL := range urls {
		if strings.Contains(articleURL, "wikipedia.org/wiki/") {
			_, err = b.SendMessage(msg.Chat.ID, articleURL, nil)
			if err != nil {
				log.Printf("Error while sending message. Err: %v\n", err)
				return
			}
			return
		}
	}

	_, err = b.SendMessage(msg.Chat.ID, "aradığın referansı bulamadım 🙈", opts)
	if err != nil {
		log.Printf("Error while sending message. Err: %v\n", err)
		return
	}
}
