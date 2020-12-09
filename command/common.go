package command

import (
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"

	customsearch "google.golang.org/api/customsearch/v1"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/googleapi/transport"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

var httpclient = &http.Client{Timeout: 10 * time.Second}

var errSearchQuotaExceeded = errors.New("Daily Limit Exceeded")

// search does a Google search with the given terms. searchType could be
// "image" or empty string.
func search(apikey, searchEngineID string, searchType string, terms ...string) ([]string, error) {
	if len(terms) == 0 {
		return nil, fmt.Errorf("no search term given")
	}

	keyword := strings.Join(terms, "+")

	imageHTTPClient := &http.Client{Transport: &transport.APIKey{Key: apikey}}
	service, err := customsearch.New(imageHTTPClient)
	if err != nil {
		return nil, fmt.Errorf("Error creating customsearch client: %v", err)
	}
	cse := customsearch.NewCseService(service)

	const imageCount = 3
	call := cse.List().Q(keyword).Cx(searchEngineID).Num(imageCount)
	if searchType == "image" {
		call = call.SearchType(searchType)
	}

	resp, err := call.Do()
	if err != nil {
		concreteErr := err.(*googleapi.Error)
		if concreteErr.Code == 403 && concreteErr.Message == "Daily Limit Exceeded" {
			return nil, errSearchQuotaExceeded
		}
		return nil, fmt.Errorf("Error making image search API call for the given criteria: %v Err: %v", keyword, err)
	}
	if len(resp.Items) == 0 {
		return nil, fmt.Errorf("Could not find any image based for the given criteria: %v", keyword)
	}

	var urls []string
	for _, url := range resp.Items {
		urls = append(urls, url.Link)
	}
	return urls, nil
}

// randChoice randomly choice an element from given elems.
func randChoice(elems []string) string {
	return elems[rand.Intn(len(elems))]
}

// emojis
const (
	emojiShrug = `¯\_(ツ)_/¯`
)
