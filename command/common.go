package command

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

const imageSearchURL = "https://ajax.googleapis.com/ajax/services/search/images"

var httpclient = http.Client{Timeout: 10 * time.Second}

// searchImage retrives an image URL for given terms.
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

// randChoice randomly choice an element from given elems
func randChoice(elems []string) string {
	return elems[rand.Intn(len(elems))]
}
