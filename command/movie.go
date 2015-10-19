package command

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

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
	// the best search engine is still google.
	// i've tried imdb, themoviedb, rottentomatoes, omdbapi.
	// themoviedb search engine was the most accurate yet still can't find any
	// result if any release date is given in query terms.
	movieAPIURL = "https://ajax.googleapis.com/ajax/services/search/web"
)

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

	arg := strings.Join(args, "+")

	u, _ := url.Parse(movieAPIURL)
	v := u.Query()
	v.Set("v", "1.0")
	v.Set("q", arg+"+movie")
	u.RawQuery = v.Encode()

	resp, err := http.Get(u.String())
	if err != nil {
		log.Printf("(movie) Error while fetching movie with given criteria: %v\n", args)
		return
	}
	defer resp.Body.Close()

	var response struct {
		ResponseData struct {
			Results []struct {
				URL   string `json:"url"`
				Title string `json:"titleNoFormatting"`
			}
		}
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		log.Printf("(movie) Error while decoding movie response: %v\n", err)
		return
	}

	for _, movie := range response.ResponseData.Results {
		if strings.Contains(movie.URL, "imdb.com/title/tt") {
			title := strings.TrimSuffix(movie.Title, " - IMDb")
			r := fmt.Sprintf("[%v](%v)", title, movie.URL)
			b.SendMessage(msg.Chat, r, tlbot.ModeMarkdown, true, nil)
			return
		}
	}

	b.SendMessage(msg.Chat, "aradığın filmi bulamadım 🙈", tlbot.ModeMarkdown, true, nil)
}

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
