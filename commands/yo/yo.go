package yo

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/igungor/ilberbot"
)

func init() {
	ilberbot.RegisterCommand("/yo", yo)
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

func yo(args ...string) string {
	if args == nil {
		log.Println("image: no argument supplied")
		return "hangi karikaturu arıyosun? ör: /yo renk dans"
	}

	arg := strings.Join(args, "+")

	u, _ := url.Parse(baseURL)
	v := u.Query()
	v.Set("q", "Yiğit+Özgür+"+arg)
	v.Set("v", "1.0")
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
		return "böyle bişey yok"
	}

	return results[0].UnescapedURL

}
