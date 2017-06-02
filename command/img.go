package command

import (
	"context"
	"fmt"
	"log"

	"github.com/igungor/ilber/bot"
	"github.com/igungor/telegram"
)

func init() {
	register(cmdImg)
}

var cmdImg = &Command{
	Name:      "img",
	ShortLine: "resim filan ara",
	Run:       runImg,
}

func runImg(ctx context.Context, b *bot.Bot, msg *telegram.Message) {
	args := msg.Args()
	md := telegram.WithParseMode(telegram.ModeMarkdown)
	if len(args) == 0 {
		term := randChoice(imgExamples)
		txt := fmt.Sprintf("ne resmi aramak istiyorsun? örneğin: */img %s*", term)
		_, err := b.SendMessage(msg.Chat.ID, txt, md)
		if err != nil {
			log.Printf("Error while sending message: %v\n", err)
		}
		return
	}

	urls, err := search(b.Config.GoogleAPIKey, b.Config.GoogleSearchEngineID, "image", args...)
	if err != nil {
		log.Printf("Error while searching image. Err: %v\n", err)
		if err == errSearchQuotaExceeded {
			_, _ = b.SendMessage(msg.Chat.ID, emojiShrug)
		}
		return
	}

	photo := telegram.Photo{
		File: telegram.File{
			URL: urls[0],
		},
	}

	_, err = b.SendPhoto(msg.Chat.ID, photo)
	if err != nil {
		log.Printf("Error while sending photo: %v\n", err)
		return
	}
}

var imgExamples = []string{
	"burdur",
	"kapadokya",
}
