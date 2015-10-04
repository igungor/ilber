package command

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/igungor/tlbot"
	"google.golang.org/api/googleapi/transport"
	youtube "google.golang.org/api/youtube/v3"
)

func init() {
	// register(cmdYoutube)
}

var cmdYoutube = &Command{
	Name:      "youtube",
	ShortLine: "vidyo filanı",
	Hide:      true,
	Run:       runYoutube,
}

const youtubeApiKey = "FIXME:XXXXXXXXXXXXXXXXXXXXXXXX"

var youtubeclient = &http.Client{Transport: &transport.APIKey{Key: youtubeApiKey}}

func runYoutube(b *tlbot.Bot, msg *tlbot.Message) {
	args := msg.Args()
	query := strings.Join(args, "+")

	service, err := youtube.New(youtubeclient)
	if err != nil {
		log.Printf("(youtube) error creating new youtube client: %v", err)
		return
	}

	call := service.Search.List("id").Type("video").Q(query).MaxResults(1)
	response, err := call.Do()
	if err != nil {
		log.Printf("(youtube) Error making youtube search API call: %v", err)
	}

	if len(response.Items) == 0 {
		b.SendMessage(msg.Chat, "aradığın videoyu bulamadım 🙈", tlbot.ModeMarkdown, true, nil)
		return
	}

	video := response.Items[0]
	v := fmt.Sprintf("https://youtube.com/watch?v=%v\n", video.Id.VideoId)
	b.SendMessage(msg.Chat, v, tlbot.ModeMarkdown, true, nil)
}
