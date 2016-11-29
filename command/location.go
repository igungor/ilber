package command

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/igungor/tlbot"
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

func runLocation(ctx context.Context, b *tlbot.Bot, msg *tlbot.Message) {
	args := msg.Args()
	if len(args) == 0 {
		_, err := b.SendMessage(msg.Chat.ID, "nerenin konumunu arayayım?", nil)
		if err != nil {
			log.Printf("Error while sending message: %v\n", err)
		}
		return
	}
	u, err := url.Parse(mapBaseURL)
	if err != nil {
		log.Fatal(err)
	}
	googleAPIKey := ctx.Value("googleAPIKey").(string)
	place := strings.Join(args, " ")
	params := u.Query()
	params.Set("key", googleAPIKey)
	params.Set("query", place)
	u.RawQuery = params.Encode()

	resp, err := httpclient.Get(u.String())
	if err != nil {
		log.Printf("Error searching place '%v'. Err: %v\n", place, err)
		return
	}
	defer resp.Body.Close()

	var places placesResponse
	if err := json.NewDecoder(resp.Body).Decode(&places); err != nil {
		log.Printf("Error decoding response. Err: %v\n", err)
		return
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("Error searching place '%v'. Status: %v\n", place, places.Status)
		return
	}

	if len(places.Results) == 0 {
		_, err = b.SendMessage(msg.Chat.ID, "bulamadim", nil)
		if err != nil {
			log.Printf("Error while sending message: %v\n", err)
		}
		return
	}

	firstPlace := places.Results[0]
	location := tlbot.Location{
		Lat:  firstPlace.Geometry.Location.Lat,
		Long: firstPlace.Geometry.Location.Long,
	}
	_, err = b.SendLocation(msg.Chat.ID, location, nil)
	if err != nil {
		log.Printf("Error sending location: %v\n", err)
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
		Icon      string   `json:"icon"`
		ID        string   `json:"id"`
		Name      string   `json:"name"`
		PlaceID   string   `json:"place_id"`
		Rating    float64  `json:"rating"`
		Reference string   `json:"reference"`
		Types     []string `json:"types"`
	} `json:"results"`
	Status string `json:"status"`
}
