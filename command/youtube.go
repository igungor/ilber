package command

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/igungor/tlbot"
	"golang.org/x/net/context"
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

func runYoutube(ctx context.Context, b *tlbot.Bot, msg *tlbot.Message) {
	googleAPIKey := ctx.Value("googleAPIKey").(string)
	youtubeHTTPClient := &http.Client{Transport: &transport.APIKey{Key: googleAPIKey}}
	service, err := youtube.New(youtubeHTTPClient)
	if err != nil {
		log.Printf("Error creating new youtube client: %v", err)
		return
	}

	args := msg.Args()
	if len(args) == 0 {
		term := randChoice(youtubeExamples)
		txt := fmt.Sprintf("ne arayayım? örneğin: */youtube %s*", term)
		err := b.SendMessage(msg.Chat.ID, txt, tlbot.ModeMarkdown, false, nil)
		if err != nil {
			log.Printf("Error while sending message: %v\n", err)
		}
		return
	}

	qs := strings.Join(args, "+")
	call := service.Search.List("id").Type("video").Q(qs).MaxResults(1)
	response, err := call.Do()
	if err != nil {
		log.Printf("Error making youtube search API call: %v", err)
		return
	}

	if len(response.Items) == 0 {
		err := b.SendMessage(msg.Chat.ID, "aradığın videoyu bulamadım 🙈", tlbot.ModeMarkdown, true, nil)
		if err != nil {
			log.Printf("Error while sending message. Err: %v\n", err)
		}
		return
	}

	video := response.Items[0]
	v := fmt.Sprintf("https://youtube.com/watch?v=%v\n", video.Id.VideoId)

	err = b.SendMessage(msg.Chat.ID, v, tlbot.ModeNone, true, nil)
	if err != nil {
		log.Printf("Error while sending message. Err: %v\n", err)
	}
}

var youtubeExamples = []string{
	"savas gayet boktan biseydir",
	"sabri bey ne yapiyorsunuz",
	"vecihi geliyor",
	"yaz kızım 200 torba çimento",
}
