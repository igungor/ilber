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

// the best search engine is still google.
// i've tried imdb, themoviedb, rottentomatoes, omdbapi.
// themoviedb search engine was the most accurate yet still can't find any
// result if any release date is given in query terms.
const movieAPIURL = "https://ajax.googleapis.com/ajax/services/search/web"

func runMovie(b *tlbot.Bot, msg *tlbot.Message) {
	args := msg.Args()
	if len(args) == 0 {
		term := randChoice(movieExamples)
		txt := fmt.Sprintf("hangi filmi arıyorsun? örneğin: */imdb %s*", term)
		err := b.SendMessage(msg.Chat, txt, tlbot.ModeMarkdown, false, nil)
		if err != nil {
			log.Printf("[movie] Error while sending message: %v\n", err)
		}
		return
	}

	qs := strings.Join(args, "+")

	u, _ := url.Parse(movieAPIURL)
	params := u.Query()
	params.Set("v", "1.0")
	params.Set("q", qs+"+movie")
	u.RawQuery = params.Encode()

	resp, err := http.Get(u.String())
	if err != nil {
		log.Printf("[movie] Error while fetching movie with given criteria '%v'. Err: %v", qs, err)
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
			err := b.SendMessage(msg.Chat, r, tlbot.ModeMarkdown, true, nil)
			if err != nil {
				log.Printf("[movie] Error while sending message. Err: %v\n", err)
			}
			return
		}
	}

	err = b.SendMessage(msg.Chat, "aradığın filmi bulamadım 🙈", tlbot.ModeMarkdown, true, nil)
	if err != nil {
		log.Printf("[movie] Error while sending message. Err: %v\n", err)
		return
	}
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
