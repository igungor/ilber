package command

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/igungor/ilber/bot"
	"github.com/igungor/telegram"
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

func runYoutube(ctx context.Context, b *bot.Bot, msg *telegram.Message) {
	youtubeHTTPClient := &http.Client{
		Transport: &transport.APIKey{
			Key: b.Config.GoogleAPIKey,
		},
	}
	service, err := youtube.New(youtubeHTTPClient)
	if err != nil {
		b.Logger.Printf("Error creating new youtube client: %v", err)
		return
	}

	args := msg.Args()
	md := telegram.WithParseMode(telegram.ModeMarkdown)
	if len(args) == 0 {
		term := randChoice(youtubeExamples)
		txt := fmt.Sprintf("ne arayayım? örneğin: */youtube %s*", term)
		_, err := b.SendMessage(msg.Chat.ID, txt, md)
		if err != nil {
			b.Logger.Printf("Error while sending message: %v\n", err)
		}
		return
	}

	qs := strings.Join(args, "+")
	call := service.Search.List("id").Type("video").Q(qs).MaxResults(1)
	response, err := call.Do()
	if err != nil {
		b.Logger.Printf("Error making youtube search API call: %v", err)
		return
	}

	if len(response.Items) == 0 {
		_, err := b.SendMessage(msg.Chat.ID, "aradığın videoyu bulamadım 🙈")
		if err != nil {
			b.Logger.Printf("Error while sending message. Err: %v\n", err)
		}
		return
	}

	video := response.Items[0]
	v := fmt.Sprintf("https://youtube.com/watch?v=%v\n", video.Id.VideoId)

	_, err = b.SendMessage(msg.Chat.ID, v)
	if err != nil {
		b.Logger.Printf("Error while sending message. Err: %v\n", err)
	}
}

var youtubeExamples = []string{
	"savas gayet boktan biseydir",
	"sabri bey ne yapiyorsunuz",
	"vecihi geliyor",
	"yaz kızım 200 torba çimento",
}
