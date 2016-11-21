package command

import (
	"fmt"
	"log"
	"strings"

	"github.com/igungor/tlbot"
	"golang.org/x/net/context"
)

func init() {
	register(cmdMovie)
}

var cmdMovie = &Command{
	Name:      "imdb",
	ShortLine: "ay-em-dii-bii",
	Run:       runMovie,
}

func runMovie(ctx context.Context, b *tlbot.Bot, msg *tlbot.Message) {
	args := msg.Args()
	if len(args) == 0 {
		term := randChoice(movieExamples)
		txt := fmt.Sprintf("hangi filmi arıyorsun? örneğin: */imdb %s*", term)
		err := b.SendMessage(msg.Chat.ID, txt, tlbot.ModeMarkdown, false, nil)
		if err != nil {
			log.Printf("Error while sending message: %v\n", err)
		}
		return
	}

	googleAPIKey := ctx.Value("googleAPIKey").(string)
	searchEngineID := ctx.Value("googleSearchEngineID").(string)

	// the best search engine is still google.
	// i've tried imdb, themoviedb, rottentomatoes, omdbapi.
	// themoviedb search engine was the most accurate yet still can't find any
	// result if any release date is given in query terms.
	urls, err := search(googleAPIKey, searchEngineID, "", args...)
	if err != nil {
		log.Printf("Error searching %v: %v\n", args, err)
		if err == errSearchQuotaExceeded {
			b.SendMessage(msg.Chat.ID, `¯\_(ツ)_/¯`, tlbot.ModeNone, false, nil)
		}
		return
	}

	for _, url := range urls {
		if strings.Contains(url, "imdb.com/title/tt") {

			err := b.SendMessage(msg.Chat.ID, url, tlbot.ModeNone, true, nil)
			if err != nil {
				log.Printf("Error while sending message. Err: %v\n", err)
			}
			return
		}
	}

	err = b.SendMessage(msg.Chat.ID, "aradığın filmi bulamadım 🙈", tlbot.ModeMarkdown, true, nil)
	if err != nil {
		log.Printf("Error while sending message. Err: %v\n", err)
		return
	}
}

var movieExamples = []string{
	"Spirited Away",
	"Mulholland Dr",
	"Oldboy",
	"Interstellar",
	"12 Angry Men",
	"Cidade de Deus",
	"The Big Lebowski",
	"There Will Be Blood",
	"Ghost in the Shell",
	"The Grey",
	"Seven Samurai",
}
