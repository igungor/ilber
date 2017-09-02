package command

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/igungor/ilber/bot"
	"github.com/igungor/telegram"
)

func init() {
	register(cmdLocation)
}

var cmdLocation = &Command{
	Name:      "konum",
	ShortLine: "ne nerde",
	Run:       runLocation,
}

const mapBaseURL = "https://maps.googleapis.com/maps/api/place/textsearch/json"

func runLocation(ctx context.Context, b *bot.Bot, msg *telegram.Message) {
	args := msg.Args()
	if len(args) == 0 {
		_, err := b.SendMessage(msg.Chat.ID, "nerenin konumunu arayayım?")
		if err != nil {
			b.Logger.Printf("Error while sending message: %v\n", err)
		}
		return
	}
	u, err := url.Parse(mapBaseURL)
	if err != nil {
		log.Fatal(err)
	}
	place := strings.Join(args, " ")
	params := u.Query()
	params.Set("key", b.Config.GoogleAPIKey)
	params.Set("query", place)
	params.Set("language", "tr")
	u.RawQuery = params.Encode()

	resp, err := httpclient.Get(u.String())
	if err != nil {
		b.Logger.Printf("Error searching place '%v'. Err: %v\n", place, err)
		return
	}
	defer resp.Body.Close()

	var places placesResponse
	if err := json.NewDecoder(resp.Body).Decode(&places); err != nil {
		b.Logger.Printf("Error decoding response. Err: %v\n", err)
		return
	}

	if resp.StatusCode != http.StatusOK {
		b.Logger.Printf("Error searching place '%v'. Status: %v\n", place, places.Status)
		return
	}

	// possible status' are at: https://developers.google.com/places/web-service/search#PlaceSearchStatusCodes
	if places.Status != "OK" {
		b.Logger.Printf("Google places query status is not OK: %v\n", places.Status)
		b.SendMessage(msg.Chat.ID, "bulamadım")
		return
	}

	if len(places.Results) == 0 {
		b.Logger.Printf("Google places query returned 0 result\n")
		b.SendMessage(msg.Chat.ID, "bulamadım")
		return
	}

	p := places.Results[0]
	venue := telegram.Venue{
		Title: fmt.Sprintf("%v (%v/5.0)", p.Name, p.Rating),
		Location: telegram.Location{
			Lat:  p.Geometry.Location.Lat,
			Long: p.Geometry.Location.Long,
		},
		Address: p.FormattedAddress,
	}

	_, err = b.SendVenue(msg.Chat.ID, venue)
	if err != nil {
		b.Logger.Printf("Error sending venue: %v\n", err)
	}
}

type placesResponse struct {
	Results []struct {
		FormattedAddress string `json:"formatted_address"`
		Geometry         struct {
			Location struct {
				Lat  float64 `json:"lat"`
				Long float64 `json:"lng"`
			} `json:"location"`
		} `json:"geometry"`
		Icon         string `json:"icon"`
		ID           string `json:"id"`
		Name         string `json:"name"`
		OpeningHours *struct {
			OpenNow bool `json:"open_now"`
		} `json:"opening_hours"`
		PlaceID           string   `json:"place_id"`
		PriceLevel        int      `json:"price_level"`
		Rating            float64  `json:"rating"`
		Reference         string   `json:"reference"`
		Types             []string `json:"types"`
		PermanentlyClosed bool     `json:"permanently_closed"`
	} `json:"results"`
	Status string `json:"status"`
}
