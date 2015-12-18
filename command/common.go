package command

import (
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	customsearch "google.golang.org/api/customsearch/v1"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/googleapi/transport"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

var (
	googleAPIKey      = os.Getenv("ILBER_GOOGLE_APIKEY")
	searchEngineID    = os.Getenv("ILBER_SEARCHENGINE_ID")
	imageclient       = &http.Client{Transport: &transport.APIKey{Key: googleAPIKey}}
	validImageFormats = []string{"png", "jpg"}

	httpclient = &http.Client{Timeout: 10 * time.Second}
)

var errImageSearchQuotaExceeded = errors.New("Daily Limit Exceeded")

// searchImage retrives an image URL for given terms.
func searchImage(terms ...string) (string, error) {
	if len(terms) == 0 {
		return "", fmt.Errorf("no search term given")
	}

	keyword := strings.Join(terms, "+")

	service, err := customsearch.New(imageclient)
	if err != nil {
		return "", fmt.Errorf("Error creating customsearch client: %v", err)
	}
	cse := customsearch.NewCseService(service)
	call := cse.List(keyword).Cx(searchEngineID).SearchType("image").Num(3)
	resp, err := call.Do()
	if err != nil {
		concreteErr := err.(*googleapi.Error)
		if concreteErr.Code == 403 && concreteErr.Message == "Daily Limit Exceeded" {
			return "", errImageSearchQuotaExceeded
		}
		return "", fmt.Errorf("Error making image search API call for the given criteria: %v Err: %v", keyword, err)
	}
	if len(resp.Items) == 0 {
		return "", fmt.Errorf("Could not find any image based for the given criteria: %v", keyword)
	}

	imageurl := resp.Items[0].Link
	return imageurl, nil
}

// randChoice randomly choice an element from given elems.
func randChoice(elems []string) string {
	return elems[rand.Intn(len(elems))]
}
