package command

import (
	"context"
	"strings"

	"github.com/igungor/ilber/bot"
	"github.com/igungor/telegram"
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

func runWiki(ctx context.Context, b *bot.Bot, msg *telegram.Message) {
	args := msg.Args()
	md := telegram.WithParseMode(telegram.ModeMarkdown)
	if len(args) == 0 {
		txt := "neye referans vereyim? mesela bana bakın: */bkz İlber Ortaylı*"
		_, err := b.SendMessage(msg.Chat.ID, txt, md)
		if err != nil {
			b.Logger.Printf("Error while sending message. Err: %v\n", err)
		}
		return
	}

	terms := []string{"wikipedia"}
	terms = append(terms, args...)

	urls, err := search(b.Config.GoogleAPIKey, b.Config.GoogleSearchEngineID, "", terms...)
	if err != nil {
		b.Logger.Printf("Error while 'bkz' query. Err: %v\n", err)
		if err == errSearchQuotaExceeded {
			b.SendMessage(msg.Chat.ID, emojiShrug)
		}
		return
	}

	for _, articleURL := range urls {
		if strings.Contains(articleURL, "wikipedia.org/wiki/") {
			_, err = b.SendMessage(msg.Chat.ID, articleURL)
			if err != nil {
				b.Logger.Printf("Error while sending message. Err: %v\n", err)
				return
			}
			return
		}
	}

	_, err = b.SendMessage(msg.Chat.ID, "aradığın referansı bulamadım 🙈")
	if err != nil {
		b.Logger.Printf("Error while sending message. Err: %v\n", err)
		return
	}
}
