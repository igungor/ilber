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
	ilberbot.RegisterCommand("/yo", image)
}

var (
	baseURL = "https://ajax.googleapis.com/ajax/services/search/images?v=1.0"
)

type Response struct {
	ResponseData struct {
		Results []Image `json:"results"`
	} `json:"responseData"`
}

type Image struct {
	UnescapedURL string `json:"unescapedURL"`
}

func image(args ...string) string {
	if args == nil {
		log.Println("image: no argument supplied")
		return ""
	}

	arg := strings.Join(args, "+")

	u, _ := url.Parse(baseURL)
	v := u.Query()
	v.Set("q", "Yiğit+Özgür+"+arg)
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

	if len(response.ResponseData.Results) == 0 {
		return "yok boyle bisi"
	}

	return response.ResponseData.Results[0].UnescapedURL

}
