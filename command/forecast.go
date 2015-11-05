package command

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/igungor/tlbot"
)

func init() {
	register(cmdForecast)
}

var cmdForecast = &Command{
	Name:      "hava",
	ShortLine: "o değil de nem fena",
	Run:       runForecast,
}

const defaultCity = "Istanbul"

var (
	apikey      = os.Getenv("ILBER_OPENWEATHERMAP_APPID")
	forecastURL = "http://api.openweathermap.org/data/2.5/weather"
)

func runForecast(b *tlbot.Bot, msg *tlbot.Message) {
	args := msg.Args()
	var location string
	if len(args) == 0 {
		location = defaultCity
	} else {
		location = strings.Join(args, " ")
	}

	u, err := url.Parse(forecastURL)
	if err != nil {
		log.Printf("[forecast] Error while parsing URL '%v'. Err: %v", forecastURL, err)
		return
	}
	params := u.Query()
	params.Set("units", "metric")
	params.Set("APPID", apikey)
	params.Set("q", location)
	u.RawQuery = params.Encode()

	resp, err := http.Get(u.String())
	if err != nil {
		log.Printf("[forecast] Error while fetching forecast for location '%v'. Err: %v\n", location, err)
		return
	}
	defer resp.Body.Close()

	var forecast Forecast
	if err := json.NewDecoder(resp.Body).Decode(&forecast); err != nil {
		log.Printf("[forecast] Error while decoding response: %v\n", err)
		return
	}

	txt := forecast.String()
	if txt == "" {
		txt = fmt.Sprintf("%v bulunamadı.", location)
	}

	err = b.SendMessage(msg.Chat.ID, txt, tlbot.ModeMarkdown, false, nil)
	if err != nil {
		log.Printf("[forecast] Error while sending message. Err: %v\n", err)
		return
	}
}

// openweathermap response
type Forecast struct {
	City    string `json:"name"`
	Weather []struct {
		ID          int    `json:"id"`
		Status      string `json:"main"`
		Description string
	}
	Temperature struct {
		Celsius float64 `json:"temp"`
	} `json:"main"`
}

func (f Forecast) String() string {
	var icon string
	now := time.Now()

	if len(f.Weather) == 0 {
		return ""
	}

	switch f.Weather[0].Status {
	case "Clear":
		if 6 < now.Hour() && now.Hour() < 18 { // for istanbul
			icon = "☀"
		} else {
			icon = "☽"
		}
	case "Clouds":
		icon = "☁"
	case "Rain":
		icon = "☔"
	case "Fog":
		icon = "▒"
	case "Mist":
		icon = "░"
	case "Haze":
		icon = "░"
	case "Snow":
		icon = "❄"
	case "Thunderstorm":
		icon = "⚡"
	default:
		icon = ""
	}

	return fmt.Sprintf("%v %v *%.1f* °C", icon, f.City, f.Temperature.Celsius)
}
