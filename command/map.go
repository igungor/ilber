package command

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/igungor/tlbot"
)

func init() {
	register(cmdMap)
}

var cmdMap = &Command{
	Name:      "konum",
	ShortLine: "ne nerde",
	Run:       runMap,
}

const mapBaseURL = "https://maps.googleapis.com/maps/api/place/textsearch/json"

func runMap(b *tlbot.Bot, msg *tlbot.Message) {
	args := msg.Args()
	if len(args) == 0 {
		if err := b.SendMessage(msg.Chat.ID, "nerenin konumunu arayayım?", tlbot.ModeNone, false, nil); err != nil {
			log.Printf("Error while sending message: %v\n", err)
		}
		return
	}
	u, err := url.Parse(mapBaseURL)
	if err != nil {
		log.Fatal(err)
	}
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
		if err := b.SendMessage(msg.Chat.ID, "bulamadim", tlbot.ModeNone, false, nil); err != nil {
			log.Printf("Error while sending message: %v\n", err)
		}
		return
	}

	firstPlace := places.Results[0]
	location := tlbot.Location{firstPlace.Geometry.Location.Lat, firstPlace.Geometry.Location.Long}
	if err := b.SendLocation(msg.Chat.ID, location, nil); err != nil {
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
