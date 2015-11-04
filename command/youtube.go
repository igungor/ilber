package command

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/igungor/tlbot"
	"google.golang.org/api/googleapi/transport"
	youtube "google.golang.org/api/youtube/v3"
)

func init() {
	register(cmdYoutube)
}

var cmdYoutube = &Command{
	Name:      "youtube",
	ShortLine: "vidyo filanı",
	Hidden:    false,
	Run:       runYoutube,
}

var (
	youtubeApiKey = os.Getenv("ILBER_YOUTUBE_APIKEY")
	youtubeclient = &http.Client{Transport: &transport.APIKey{Key: youtubeApiKey}}
)

func runYoutube(b *tlbot.Bot, msg *tlbot.Message) {
	service, err := youtube.New(youtubeclient)
	if err != nil {
		log.Printf("[youtube] error creating new youtube client: %v", err)
		return
	}

	args := msg.Args()
	if len(args) == 0 {
		term := randChoice(youtubeExamples)
		txt := fmt.Sprintf("ne arayayım? örneğin: */youtube %s*", term)
		err := b.SendMessage(msg.Chat, txt, tlbot.ModeMarkdown, false, nil)
		if err != nil {
			log.Printf("[movie] Error while sending message: %v\n", err)
		}
		return
	}

	qs := strings.Join(args, "+")
	call := service.Search.List("id").Type("video").Q(qs).MaxResults(1)
	response, err := call.Do()
	if err != nil {
		log.Printf("[youtube] Error making youtube search API call: %v", err)
	}

	if len(response.Items) == 0 {
		err := b.SendMessage(msg.Chat, "aradığın videoyu bulamadım 🙈", tlbot.ModeMarkdown, true, nil)
		if err != nil {
			log.Printf("[youtube] Error while sending message. Err: %v\n", err)
		}
		return
	}

	video := response.Items[0]
	v := fmt.Sprintf("https://youtube.com/watch?v=%v\n", video.Id.VideoId)

	err = b.SendMessage(msg.Chat, v, tlbot.ModeNone, true, nil)
	if err != nil {
		log.Printf("[youtube] Error while sending message. Err: %v\n", err)
	}
}

var youtubeExamples = []string{
	"savas gayet boktan biseydir",
	"sabri bey ne yapiyorsunuz",
	"vecihi geliyor",
	"yaz kızım 200 torba çimento",
}
