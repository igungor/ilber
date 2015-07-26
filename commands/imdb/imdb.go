package imdb

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/igungor/ilberbot"
)

var (
	// Use upstream api if needed: http://www.imdb.com/xml/find?json=1&nr=1&tt=on&q=lost
	baseURL      = "http://www.omdbapi.com/"
	imdbTitleURL = "http://www.imdb.com/title/"
	client       = http.Client{Timeout: 10 * time.Second}
)

func init() {
	ilberbot.RegisterCommand("/imdb", imdb)
}

type Response struct {
	ID       string `json:"imdbID"`
	Rating   string `json:"imdbRating"`
	Type     string
	Title    string
	Year     string
	Released string
	Runtime  string
	Genre    string
	Director string
	Writer   string
	Actors   string
	Plot     string
	Language string
	Country  string
}

func imdb(args ...string) string {
	if args == nil {
		log.Println("image: no argument supplied")
		return ""
	}

	arg := strings.Join(args, "+")

	u, _ := url.Parse(baseURL)
	v := u.Query()
	v.Set("t", arg)    // movie title
	v.Set("r", "json") // return type
	u.RawQuery = v.Encode()

	resp, err := http.Get(u.String())
	if err != nil {
		log.Println(err)
		return ""
	}
	defer resp.Body.Close()

	var response Response
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		log.Printf("decode error: %v", err)
		return ""
	}

	if response.ID == "" {
		log.Printf("imdb: no title found")
		return ""
	}

	return imdbTitleURL + response.ID
}
