package command

import (
	"fmt"
	"log"

	"github.com/igungor/tlbot"
	"golang.org/x/net/context"
)

func init() {
	register(cmdYo)
}

var cmdYo = &Command{
	Name:      "yo",
	ShortLine: "yiğit özgür şeysi",
	Run:       runYo,
}

func runYo(ctx context.Context, b *tlbot.Bot, msg *tlbot.Message) {
	args := msg.Args()
	if len(args) == 0 {
		term := randChoice(yoExamples)
		txt := fmt.Sprintf("hangi karikatürü arıyorsun? örneğin: */yo %s*", term)
		err := b.SendMessage(msg.Chat.ID, txt, tlbot.ModeMarkdown, false, nil)
		if err != nil {
			log.Printf("Error while sending message: %v\n", err)
		}
		return
	}

	googleAPIKey := ctx.Value("googleAPIKey").(string)
	searchEngineID := ctx.Value("googleSearchEngineID").(string)

	terms := []string{"Yiğit", "Özgür"}
	terms = append(terms, args...)
	u, err := search(googleAPIKey, searchEngineID, "image", terms...)
	if err != nil {
		log.Printf("Error while searching image with given criteria: %v. Err: %v\n", args, err)
		if err == errSearchQuotaExceeded {
			b.SendMessage(msg.Chat.ID, `¯\_(ツ)_/¯`, tlbot.ModeNone, false, nil)
		}
		return
	}

	photo := tlbot.Photo{File: tlbot.File{FileURL: u[0]}}
	err = b.SendPhoto(msg.Chat.ID, photo, "", nil)
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
