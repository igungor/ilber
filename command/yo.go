package command

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/igungor/tlbot"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
	register(cmdYo)
}

var cmdYo = &Command{
	Name:      "yo",
	ShortLine: "yiğit özgür şeysi",
	Run:       runYo,
}

var yoExamples = []string{
	"renk dans",
	"bağa mı didin",
	"düşünemedi",
	"lütfen olsun çünkü",
	"harika adam",
	"sipirmin",
	"lanet olsun",
	"flemenko",
}

func runYo(b *tlbot.Bot, msg *tlbot.Message) {
	args := msg.Args()

	if len(args) == 0 {
		term := yoExamples[rand.Intn(len(yoExamples))]
		txt := fmt.Sprintf("hangi karikatürü arıyorsun? örneğin: */yo %s*", term)
		err := b.SendMessage(msg.From, txt, tlbot.ModeMarkdown, false, nil)
		if err != nil {
			log.Printf("(yo) Error while sending message: %v\n", err)
		}
		return
	}

	terms := []string{"Yiğit", "Özgür"}
	terms = append(terms, args...)

	u, err := searchImage(terms...)
	if err != nil {
		log.Println("(img) Error while searching image with given criteria: %v\n", args)
		return
	}

	photo := tlbot.Photo{File: tlbot.File{FileURL: u}}
	b.SendPhoto(msg.From, photo, "", nil)
}
