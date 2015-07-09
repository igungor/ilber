package caps

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strings"

	"github.com/igungor/ilberbot"
)

func init() {
	ilberbot.RegisterCommand("/caps", caps)
}

var (
	baseURL = "https://ajax.googleapis.com/ajax/services/search/images"
)

type Response struct {
	ResponseData struct {
		Results []Image `json:"results"`
	} `json:"responseData"`
}

type Image struct {
	UnescapedURL string `json:"unescapedURL"`
}

func caps(args ...string) string {
	if args == nil {
		log.Println("image: no argument supplied")
		return ""
	}

	arg := strings.Join(args, "+")

	u, _ := url.Parse(baseURL)
	v := u.Query()
	v.Set("q", "caps+"+arg) // query (mandatory)
	v.Set("v", "1.0")       // version (mandatory)
	v.Set("rsz", "3")       // result size (optional)
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

	results := response.ResponseData.Results

	if len(results) == 0 {
		return "yok boyle bisi"
	}

	return results[rand.Intn(len(results))].UnescapedURL

}
