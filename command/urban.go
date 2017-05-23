package command

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
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
		_, err := b.SendMessage(msg.Chat.ID, txt, nil)
		if err != nil {
			log.Printf("Error sending message: %v\n", err)
		}
		return
	}

	u, err := url.Parse(urbanURL)
	if err != nil {
		log.Printf("Error parsing Urban Dictionary URL: %v\n", err)
		b.SendMessage(msg.Chat.ID, emojiShrug, nil)
		return
	}

	term := strings.Join(args, " ")
	params := url.Values{}
	params.Set("term", term)
	u.RawQuery = params.Encode()

	resp, err := httpclient.Get(u.String())
	if err != nil {
		log.Printf("Error sending request to Urban Dictionary API: %v\n", err)
		b.SendMessage(msg.Chat.ID, emojiShrug, nil)
		return
	}
	defer resp.Body.Close()

	var r response
	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		log.Printf("Error parsing response body from Urban Dictonary: %v\n", err)
		b.SendMessage(msg.Chat.ID, emojiShrug, nil)
		return
	}

	if r.ResultType == "no_results" {
		log.Printf("Empty result set from Urban Dictionary for term %q\n", term)
		b.SendMessage(msg.Chat.ID, fmt.Sprintf("UrbanDictonary'de %q diye birşey yok", term), nil)
		return
	}

	_, err = b.SendMessage(msg.Chat.ID, r.String(), nil)
	if err != nil {
		log.Printf("Error sending message to Telegram: %v\n", err)
	}
}

type response struct {
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

func (r response) String() string {
	if len(r.List) == 0 {
		return fmt.Sprintf("UrbanDictonary'de böyle birşey yok")
	}

	var maxItems int
	if len(r.List) > 3 {
		maxItems = 3
	} else {
		maxItems = len(r.List)
	}

	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("Definitions of %q\n\n", r.List[0].Word))
	for i := 0; i < maxItems; i++ {
		item := r.List[i]
		buf.WriteString(fmt.Sprintf("* %v\n", item.Definition))
		buf.WriteString("\n")
	}

	return buf.String()
}
