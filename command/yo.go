package command

import (
	"fmt"
	"log"

	"github.com/igungor/tlbot"
)

func init() {
	register(cmdYo)
}

var cmdYo = &Command{
	Name:      "yo",
	ShortLine: "yiğit özgür şeysi",
	Run:       runYo,
}

func runYo(b *tlbot.Bot, msg *tlbot.Message) {
	args := msg.Args()

	if len(args) == 0 {
		term := randChoice(yoExamples)
		txt := fmt.Sprintf("hangi karikatürü arıyorsun? örneğin: */yo %s*", term)
		err := b.SendMessage(msg.Chat, txt, tlbot.ModeMarkdown, false, nil)
		if err != nil {
			log.Printf("[yo] Error while sending message: %v\n", err)
		}
		return
	}

	terms := []string{"Yiğit", "Özgür"}
	terms = append(terms, args...)

	u, err := searchImage(terms...)
	if err != nil {
		log.Printf("[yo] Error while searching image with given criteria: %v\n", args)
		return
	}

	photo := tlbot.Photo{File: tlbot.File{FileURL: u}}
	err = b.SendPhoto(msg.Chat, photo, "", nil)
	if err != nil {
		log.Printf("[yo] Error while sending image: %v\n", err)
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
