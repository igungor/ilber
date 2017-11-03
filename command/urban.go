package command

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/igungor/ilber/bot"
	"github.com/igungor/telegram"
)

func init() {
	register(cmdUrban)
}

var cmdUrban = &Command{
	Name:      "urban",
	ShortLine: "urban dictionary",
	Run:       runUrban,
	Hidden:    true,
}

const urbanURL = "http://api.urbandictionary.com/v0/define"

func runUrban(ctx context.Context, b *bot.Bot, msg *telegram.Message) {
	args := msg.Args()
	if len(args) == 0 {
		txt := "neyi arayayım?"
		_, err := b.SendMessage(msg.Chat.ID, txt)
		if err != nil {
			b.Logger.Printf("Error sending message: %v\n", err)
		}
		return
	}

	u, err := url.Parse(urbanURL)
	if err != nil {
		b.Logger.Printf("Error parsing Urban Dictionary URL: %v\n", err)
		b.SendMessage(msg.Chat.ID, emojiShrug)
		return
	}

	term := strings.Join(args, " ")
	params := url.Values{}
	params.Set("term", term)
	u.RawQuery = params.Encode()

	resp, err := httpclient.Get(u.String())
	if err != nil {
		b.Logger.Printf("Error sending request to Urban Dictionary API: %v\n", err)
		b.SendMessage(msg.Chat.ID, emojiShrug)
		return
	}
	defer resp.Body.Close()

	var r urbanResponse
	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		b.Logger.Printf("Error parsing response body from Urban Dictonary: %v\n", err)
		b.SendMessage(msg.Chat.ID, emojiShrug)
		return
	}

	if r.ResultType == "no_results" {
		b.Logger.Printf("Empty result set from Urban Dictionary for term %q\n", term)
		b.SendMessage(msg.Chat.ID, fmt.Sprintf("UrbanDictonary'de %q diye birşey yok", term))
		return
	}

	_, err = b.SendMessage(msg.Chat.ID, r.String())
	if err != nil {
		b.Logger.Printf("Error sending message to Telegram: %v\n", err)
	}
}

type urbanResponse struct {
	List []struct {
		Author      string `json:"author"`
		CurrentVote string `json:"current_vote"`
		Defid       int    `json:"defid"`
		Definition  string `json:"definition"`
		Example     string `json:"example"`
		Permalink   string `json:"permalink"`
		ThumbsDown  int    `json:"thumbs_down"`
		ThumbsUp    int    `json:"thumbs_up"`
		Word        string `json:"word"`
	} `json:"list"`
	ResultType string        `json:"result_type"`
	Sounds     []interface{} `json:"sounds"`
	Tags       []string      `json:"tags"`
}

func (ur urbanResponse) String() string {
	if len(ur.List) == 0 {
		return fmt.Sprintf("UrbanDictonary'de böyle birşey yok")
	}

	var maxItems int
	if len(ur.List) > 3 {
		maxItems = 3
	} else {
		maxItems = len(ur.List)
	}

	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("Definitions of %q\n\n", ur.List[0].Word))
	for i := 0; i < maxItems; i++ {
		item := ur.List[i]
		buf.WriteString(fmt.Sprintf("* %v\n", item.Definition))
		buf.WriteString("\n")
	}

	return buf.String()
}
