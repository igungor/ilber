package command

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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

var forecastURL = "http://api.openweathermap.org/data/2.5/weather?q=%s&units=metric"

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

func runForecast(b *tlbot.Bot, msg *tlbot.Message) {
	args := msg.Args()
	var location string
	if len(args) == 0 {
		location = "Istanbul"
	} else {
		location = strings.Join(args, " ")
	}

	url := fmt.Sprintf(forecastURL, location)
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("(yo) Error while fetching forecast for location '%v'. Err: %v\n", location, err)
		return
	}
	defer resp.Body.Close()

	var forecast Forecast
	if err := json.NewDecoder(resp.Body).Decode(&forecast); err != nil {
		log.Printf("(forecast) Error while decoding response: %v\n", err)
		return
	}

	if forecast.String() == "" {
		b.SendMessage(msg.Chat, fmt.Sprintf("%v bulunamadı.", location), tlbot.ModeNone, false, nil)
		return
	}

	b.SendMessage(msg.Chat, forecast.String(), tlbot.ModeMarkdown, false, nil)
}
