package command

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/igungor/ilber/bot"
	"github.com/igungor/telegram"
)

func init() {
	register(cmdForecast)
}

var cmdForecast = &Command{
	Name:      "hava",
	ShortLine: "o değil de nem fena",
	Run:       runForecast,
}

const (
	defaultCity = "Istanbul"
	forecastURL = "http://api.openweathermap.org/data/2.5/weather"
)

func runForecast(ctx context.Context, b *bot.Bot, msg *telegram.Message) {
	args := msg.Args()
	var location string
	if len(args) == 0 {
		location = defaultCity
	} else {
		location = strings.Join(args, " ")
	}

	u, err := url.Parse(forecastURL)
	if err != nil {
		log.Printf("Error while parsing URL '%v'. Err: %v", forecastURL, err)
		return
	}

	params := u.Query()
	params.Set("units", "metric")
	params.Set("APPID", b.Config.OpenweathermapAppID)
	params.Set("q", location)
	u.RawQuery = params.Encode()

	resp, err := http.Get(u.String())
	if err != nil {
		log.Printf("Error while fetching forecast for location '%v'. Err: %v\n", location, err)
		return
	}
	defer resp.Body.Close()

	var forecast forecast
	if err := json.NewDecoder(resp.Body).Decode(&forecast); err != nil {
		log.Printf("Error while decoding response: %v\n", err)
		return
	}

	txt := forecast.String()
	if txt == "" {
		txt = fmt.Sprintf("%v bulunamadı.", location)
	}
	_, err = b.SendMessage(msg.Chat.ID, txt, telegram.WithParseMode(telegram.ModeMarkdown))
	if err != nil {
		log.Printf("Error while sending message. Err: %v\n", err)
		return
	}
}

// openweathermap response
type forecast struct {
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

func (f forecast) String() string {
	if len(f.Weather) == 0 {
		return ""
	}

	var icon string
	now := time.Now()
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
