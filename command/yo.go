package command

import (
	"context"
	"fmt"
	"log"

	"github.com/igungor/ilber/bot"
	"github.com/igungor/telegram"
)

func init() {
	register(cmdYo)
}

var cmdYo = &Command{
	Name:      "yo",
	ShortLine: "yiğit özgür şeysi",
	Run:       runYo,
}

func runYo(ctx context.Context, b *bot.Bot, msg *telegram.Message) {
	args := msg.Args()
	md := telegram.WithParseMode(telegram.ModeMarkdown)
	if len(args) == 0 {
		term := randChoice(yoExamples)
		txt := fmt.Sprintf("hangi karikatürü arıyorsun? örneğin: */yo %s*", term)
		_, err := b.SendMessage(msg.Chat.ID, txt, md)
		if err != nil {
			log.Printf("Error while sending message: %v\n", err)
		}
		return
	}

	terms := []string{"Yiğit", "Özgür"}
	terms = append(terms, args...)
	u, err := search(b.Config.GoogleAPIKey, b.Config.GoogleSearchEngineID, "image", terms...)
	if err != nil {
		log.Printf("Error while searching image with given criteria: %v. Err: %v\n", args, err)
		if err == errSearchQuotaExceeded {
			_, _ = b.SendMessage(msg.Chat.ID, emojiShrug)
		}
		return
	}

	photo := telegram.Photo{
		File: telegram.File{
			URL: u[0],
		},
	}
	_, err = b.SendPhoto(msg.Chat.ID, photo)
	if err != nil {
		log.Printf("Error while sending image: %v\n", err)
		return
	}
}

var yoExamples = []string{
	"renk dans",
	"bağa mı didin",
	"düşünemedi",
	"lütfen olsun çünkü",
	"geldi yine",
	"sipirmin",
	"lanet olsun sana",
	"flemenko",
}
