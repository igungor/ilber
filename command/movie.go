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
	movieAPIURL  = "http://www.imdb.com/xml/find?json=1&nr=1&tt=on&q=%s"
	imdbTitleURL = "http://www.imdb.com/title/"
	maxResult    = 4
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

	arg := strings.Join(args, "+")

	u, _ := url.Parse(movieAPIURL)
	v := u.Query()
	v.Set("json", "1")
	v.Set("nr", "1")
	v.Set("tt", "on")
	v.Set("q", arg)
	u.RawQuery = v.Encode()

	resp, err := http.Get(u.String())
	if err != nil {
		log.Printf("(movie) Error while fetching movie with given criteria: %v. Err: %v\n", args, err)
		return
	}
	defer resp.Body.Close()

	var response imdbResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		log.Printf("(movie) Error while decoding movie response: %v\n", err)
		return
	}

	opts := &tlbot.SendOptions{
		ReplyMarkup: tlbot.ReplyMarkup{
			Resize:    true,
			OneTime:   true,
			Selective: true,
		},
	}

	spew.Config.DisableMethods = true
	spew.Dump(msg)

	movies := consolidateMovies(response)
	switch len(movies) {
	case 0:
		log.Printf("(movie) No title found with given criteria: %v\n", arg)
		return
	case 1:
		movie := movies[0]
		year := strings.Split(movie.Description, ",")[0]
		r := fmt.Sprintf("[%v (%v)](%v%v)", movie.Title, year, imdbTitleURL, movie.ID)
		b.SendMessage(msg.Chat, r, tlbot.ModeMarkdown, true, nil)
		return
	default:
		fillKeyboard(&opts.ReplyMarkup, movies)
		b.SendMessage(msg.Chat, "birden fazla film var. hangisi?", tlbot.ModeMarkdown, true, opts)
		return
	}
}

// imdb search results are so confusing. join the exact and popular results.
func consolidateMovies(response imdbResponse) []movie {
	var movies []movie
	switch len(response.TitleExact) {
	default:
		movies = append(movies, response.TitleExact[0])
		if len(response.TitleExact) > 1 {
			movies = append(movies, response.TitleExact[1])
		}
		fallthrough
	case 0:
		movies = append(movies, response.TitlePopular...)
	}
	return movies
}

func fillKeyboard(kbd *tlbot.ReplyMarkup, movies []movie) {
	if len(movies) > maxResult {
		movies = movies[:maxResult]
	}

	kbd.Keyboard = make([][]string, len(movies))
	for i, movie := range movies {
		year := strings.Split(movie.Description, ",")[0]
		item := movie.Title + "," + year
		kbd.Keyboard[i] = []string{item}
	}
}

type imdbResponse struct {
	TitleExact     []movie `json:"title_exact"`
	TitlePopular   []movie `json:"title_popular"`
	TitleSubstring []movie `json:"title_substring"`
	TitleApprox    []movie `json:"title_approx"`
}

type movie struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"title_description"`
}
