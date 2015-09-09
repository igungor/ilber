package command

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/igungor/tlbot"
)

func init() {
	rand.Seed(time.Now().UnixNano())
	register(cmdYo)
}

var cmdYo = &Command{
	Name: "yo",
	Run:  runYo,
}

const imageSearchURL = "https://ajax.googleapis.com/ajax/services/search/images"

var examples = []string{
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
		term := examples[rand.Intn(len(examples))]
		txt := fmt.Sprintf("hangi karikatürü arıyosun? örneğin: */yo %s*", term)
		err := b.SendMessage(msg.From, txt, tlbot.ModeMarkdown, false, nil)
		if err != nil {
			log.Printf("(yo) Error while sending message: %v\n", err)
		}
		return
	}

	keyword := strings.Join(args, "+")

	u, _ := url.Parse(imageSearchURL)
	v := u.Query()
	v.Set("q", "Yiğit+Özgür+"+keyword)
	v.Set("v", "1.0")
	u.RawQuery = v.Encode()

	resp, err := http.Get(u.String())
	if err != nil {
		log.Printf("(yo) Error while fetching image for the criteria '%v'. Err: %v\n", keyword, err)
		return
	}
	defer resp.Body.Close()

	// datastructure of image search response
	var response struct {
		ResponseData struct {
			Results []struct {
				UnescapedURL string `json:"unescapedURL"`
			} `json:"results"`
		} `json:"responseData"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		log.Printf("(yo) Error while decoding response: %v\n", err)
		return
	}

	results := response.ResponseData.Results
	if len(results) == 0 {
		b.SendMessage(msg.From, "_böyle bişey yok_", tlbot.ModeMarkdown, false, nil)
		return
	}

	photo := tlbot.Photo{File: tlbot.File{FileURL: results[0].UnescapedURL}}
	b.SendPhoto(msg.From, photo, "", nil)
}
