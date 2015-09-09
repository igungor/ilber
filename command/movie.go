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
	register(cmdMovie)
}

var cmdMovie = &Command{
	Name:      "imdb",
	ShortLine: "ay-em-dii-bii",
	Run:       runMovie,
}

const (
	// Use upstream api if needed: http://www.imdb.com/xml/find?json=1&nr=1&tt=on&q=lost
	movieAPIURL  = "http://www.omdbapi.com/"
	imdbTitleURL = "http://www.imdb.com/title/"
)

var movieExamples = []string{
	"Spirited Away",
	"Mulholland Dr",
	"Oldboy",
	"Interstellar",
	"12 Angry Men",
	"Cidade de Deus",
	"The Big Lebowski",
	"There Will Be Blood",
}

func runMovie(b *tlbot.Bot, msg *tlbot.Message) {
	args := msg.Args()

	if len(args) == 0 {
		term := movieExamples[rand.Intn(len(movieExamples))]
		txt := fmt.Sprintf("hangi filmi arıyorsun? örneğin: */imdb %s*", term)
		err := b.SendMessage(msg.From, txt, tlbot.ModeMarkdown, false, nil)
		if err != nil {
			log.Printf("(movie) Error while sending message: %v\n", err)
		}
		return
	}

	arg := strings.Join(args, "+")

	u, _ := url.Parse(movieAPIURL)
	v := u.Query()
	v.Set("t", arg)    // movie title
	v.Set("r", "json") // return type
	u.RawQuery = v.Encode()

	resp, err := http.Get(u.String())
	if err != nil {
		log.Printf("(movie) Error while fetching movie with given criteria: %v\n", args)
		return
	}
	defer resp.Body.Close()

	var response struct {
		ID     string `json:"imdbID"`
		Rating string `json:"imdbRating"`
		Title  string
		Year   string
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		log.Printf("(movie) Error while decoding movie response: %v\n", err)
		return
	}

	if response.ID == "" {
		log.Printf("(movie) No title found with given criteria: %v\n", arg)
		return
	}

	r := imdbTitleURL + response.ID

	// enable preview
	b.SendMessage(msg.From, r, tlbot.ModeNone, true, nil)
}
