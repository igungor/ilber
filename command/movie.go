package command

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/igungor/tlbot"
)

func init() {
	register(cmdMovie)
}

var cmdMovie = &Command{
	Name:      "imdb",
	ShortLine: "ay-em-dii-bii",
	Run:       runMovie,
}

const (
	// Use upstream api if needed: http://www.imdb.com/xml/find?json=1&nr=1&tt=on&q=lost
	movieAPIURL  = "https://api.themoviedb.org/3/search/movie"
	imdbTitleURL = "http://www.imdb.com/title/"
	// FIXME!!!
	ApiKey = "55a80aeef3620a31015d41be2a1c4fbc"
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
		term := randChoice(movieExamples)
		txt := fmt.Sprintf("hangi filmi arıyorsun? örneğin: */imdb %s*", term)
		err := b.SendMessage(msg.Chat, txt, tlbot.ModeMarkdown, false, nil)
		if err != nil {
			log.Printf("(movie) Error while sending message: %v\n", err)
		}
		return
	}

	arg := strings.Join(args, " ")

	u, _ := url.Parse(movieAPIURL)
	v := u.Query()
	v.Set("api_key", ApiKey)
	v.Set("query", arg)
	v.Set("year", "2015")
	u.RawQuery = v.Encode()

	fmt.Println(u.String())
	resp, err := http.Get(u.String())
	if err != nil {
		log.Printf("(movie) Error while fetching movie with given criteria: %v\n", args)
		return
	}
	defer resp.Body.Close()

	var r tmdbResponse
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		log.Printf("(movie) Error while decoding movie response: %v\n", err)
		return
	}

	spew.Dump(r)
}

type tmdbResponse struct {
	Results []struct {
		ID            int     `json:"id"`
		OriginalTitle string  `json:"original_title"`
		Popularity    float64 `json:"popularity"`
		ReleaseDate   string  `json:"release_date"`
		Title         string  `json:"title"`
		VoteAverage   float64 `json:"vote_average"`
		VoteCount     int     `json:"vote_count"`
	} `json:"results"`
}
