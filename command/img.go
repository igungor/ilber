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
	rand.Seed(time.Now().UTC().UnixNano())
	register(cmdImg)
}

const imageSearchURL = "https://ajax.googleapis.com/ajax/services/search/images"

var cmdImg = &Command{
	Name: "img",
	Run:  runImg,
}

var imgExamples = []string{
	"burdur",
	"kapadoky",
	"leblebi tozu",
}

func runImg(b *tlbot.Bot, msg *tlbot.Message) {
	args := msg.Args()

	if len(args) == 0 {
		term := imgExamples[rand.Intn(len(imgExamples))]
		txt := fmt.Sprintf("ne resmi aramak istiyorsun? örneğin: */img %s*", term)
		err := b.SendMessage(msg.From, txt, tlbot.ModeMarkdown, false, nil)
		if err != nil {
			log.Printf("(img) Error while sending message: %v\n", err)
		}
		return
	}

	u, err := searchImage(args...)
	if err != nil {
		log.Println("(img) Error while searching image with given criteria: %v\n", args)
		return
	}

	photo := tlbot.Photo{File: tlbot.File{FileURL: u}}
	b.SendPhoto(msg.From, photo, "", nil)
}

func searchImage(terms ...string) (string, error) {
	if len(terms) == 0 {
		return "", fmt.Errorf("no search term given")
	}

	keyword := strings.Join(terms, "+")

	u, _ := url.Parse(imageSearchURL)
	v := u.Query()
	v.Set("q", keyword)
	v.Set("v", "1.0")
	u.RawQuery = v.Encode()

	resp, err := http.Get(u.String())
	if err != nil {
		return "", err
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
		return "", err
	}

	results := response.ResponseData.Results
	if len(results) == 0 {
		return "", fmt.Errorf("no results for the given criteria: %v\n", keyword)
	}

	return results[0].UnescapedURL, nil
}
