package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

var (
	dakickToken = os.Getenv("DAKICK_TOKEN")
	dakickURL   = "https://api.dakick.com/api/v1/movies/in-theaters?location_ids=207&per_page=15"
)

func init() {
	register("/vizyon", theaters)
}

// dakick response
type (
	Movie struct {
		Results []Result
	}

	Result struct {
		Name string
		AKA  string `json:"aka"`
	}
)

func (m Movie) Filter(criteria string) Movie {
	var results []Result
	for _, movie := range m.Results {
		if strings.Contains(movie.Name, criteria) {
			continue
		}
		results = append(results, movie)
	}

	return Movie{results}
}

func (m Movie) String() string {
	var buf bytes.Buffer

	buf.WriteString("🎦 Istanbul'da vizyon filmleri\n")
	for _, movie := range m.Results {
		if len(movie.Name) > 30 {
			movie.Name = movie.Name[:30] + "…"
		}

		if movie.AKA != "" {
			buf.WriteString(fmt.Sprintf("- %v (%v)\n", movie.AKA, movie.Name))
		} else {
			buf.WriteString(fmt.Sprintf("- %v\n", movie.Name))
		}
	}

	return buf.String()
}

func theaters(args ...string) string {
	if dakickToken == "" {
		log.Println("DAKICK_TOKEN must be set")
		return ""
	}

	var client http.Client
	req, err := http.NewRequest("GET", dakickURL, nil)
	if err != nil {
		log.Printf("theaters request error: %v\n", err)
		return ""
	}
	req.Header.Set("X-DAKICK-API-TOKEN", dakickToken)

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("theaters client.do error: %v\n", err)
		return ""
	}
	defer resp.Body.Close()

	var movie Movie
	if err := json.NewDecoder(resp.Body).Decode(&movie); err != nil {
		log.Printf("decode error: %v", err)
		return ""
	}

	movies := movie.Filter("3D")
	printdebug("%v\n", movies)

	return movies.String()
}
