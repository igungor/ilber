package command

import (
	"context"
	"encoding/json"
	"fmt"
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
		b.Logger.Printf("Error while parsing URL '%v'. Err: %v", forecastURL, err)
		return
	}

	params := u.Query()
	params.Set("units", "metric")
	params.Set("APPID", b.Config.OpenweathermapAppID)
	params.Set("q", location)
	u.RawQuery = params.Encode()

	resp, err := http.Get(u.String())
	if err != nil {
		b.Logger.Printf("Error while fetching forecast for location '%v'. Err: %v\n", location, err)
		return
	}
	defer resp.Body.Close()

	var forecast forecast
	if err := json.NewDecoder(resp.Body).Decode(&forecast); err != nil {
		b.Logger.Printf("Error while decoding response: %v\n", err)
		return
	}

	txt := forecast.String()
	if txt == "" {
		txt = fmt.Sprintf("%v bulunamadı.", location)
	}
	_, err = b.SendMessage(msg.Chat.ID, txt, telegram.WithParseMode(telegram.ModeMarkdown))
	if err != nil {
		b.Logger.Printf("Error while sending message. Err: %v\n", err)
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
	Sys struct {
		Country string `json:"country"`
		Sunrise int64  `json:"sunrise"`
		Sunset  int64  `json:"sunset"`
	} `json:"sys"`
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
		sunrise := time.Unix(f.Sys.Sunrise, 0)
		sunset := time.Unix(f.Sys.Sunset, 0)
		if now.After(sunrise) && now.Before(sunset) {
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
	return fmt.Sprintf("%v %v%v *%.1f* °C", icon, f.City, flagLookup(f.Sys.Country), f.Temperature.Celsius)
}

func flagLookup(code string) string {
	var flags = map[string]string{
		"AD": "🇦🇩",
		"AE": "🇦🇪",
		"AF": "🇦🇫",
		"AG": "🇦🇬",
		"AI": "🇦🇮",
		"AL": "🇦🇱",
		"AM": "🇦🇲",
		"AO": "🇦🇴",
		"AQ": "🇦🇶",
		"AR": "🇦🇷",
		"AS": "🇦🇸",
		"AT": "🇦🇹",
		"AU": "🇦🇺",
		"AW": "🇦🇼",
		"AX": "🇦🇽",
		"AZ": "🇦🇿",
		"BA": "🇧🇦",
		"BB": "🇧🇧",
		"BD": "🇧🇩",
		"BE": "🇧🇪",
		"BF": "🇧🇫",
		"BG": "🇧🇬",
		"BH": "🇧🇭",
		"BI": "🇧🇮",
		"BJ": "🇧🇯",
		"BL": "🇧🇱",
		"BM": "🇧🇲",
		"BN": "🇧🇳",
		"BO": "🇧🇴",
		"BQ": "🇧🇶",
		"BR": "🇧🇷",
		"BS": "🇧🇸",
		"BT": "🇧🇹",
		"BV": "🇧🇻",
		"BW": "🇧🇼",
		"BY": "🇧🇾",
		"BZ": "🇧🇿",
		"CA": "🇨🇦",
		"CC": "🇨🇨",
		"CD": "🇨🇩",
		"CF": "🇨🇫",
		"CG": "🇨🇬",
		"CH": "🇨🇭",
		"CI": "🇨🇮",
		"CK": "🇨🇰",
		"CL": "🇨🇱",
		"CM": "🇨🇲",
		"CN": "🇨🇳",
		"CO": "🇨🇴",
		"CR": "🇨🇷",
		"CU": "🇨🇺",
		"CV": "🇨🇻",
		"CW": "🇨🇼",
		"CX": "🇨🇽",
		"CY": "🇨🇾",
		"CZ": "🇨🇿",
		"DE": "🇩🇪",
		"DJ": "🇩🇯",
		"DK": "🇩🇰",
		"DM": "🇩🇲",
		"DO": "🇩🇴",
		"DZ": "🇩🇿",
		"EC": "🇪🇨",
		"EE": "🇪🇪",
		"EG": "🇪🇬",
		"EH": "🇪🇭",
		"ER": "🇪🇷",
		"ES": "🇪🇸",
		"ET": "🇪🇹",
		"EU": "🇪🇺",
		"FI": "🇫🇮",
		"FJ": "🇫🇯",
		"FK": "🇫🇰",
		"FM": "🇫🇲",
		"FO": "🇫🇴",
		"FR": "🇫🇷",
		"GA": "🇬🇦",
		"GB": "🇬🇧",
		"GD": "🇬🇩",
		"GE": "🇬🇪",
		"GF": "🇬🇫",
		"GG": "🇬🇬",
		"GH": "🇬🇭",
		"GI": "🇬🇮",
		"GL": "🇬🇱",
		"GM": "🇬🇲",
		"GN": "🇬🇳",
		"GP": "🇬🇵",
		"GQ": "🇬🇶",
		"GR": "🇬🇷",
		"GS": "🇬🇸",
		"GT": "🇬🇹",
		"GU": "🇬🇺",
		"GW": "🇬🇼",
		"GY": "🇬🇾",
		"HK": "🇭🇰",
		"HM": "🇭🇲",
		"HN": "🇭🇳",
		"HR": "🇭🇷",
		"HT": "🇭🇹",
		"HU": "🇭🇺",
		"ID": "🇮🇩",
		"IE": "🇮🇪",
		"IL": "🇮🇱",
		"IM": "🇮🇲",
		"IN": "🇮🇳",
		"IO": "🇮🇴",
		"IQ": "🇮🇶",
		"IR": "🇮🇷",
		"IS": "🇮🇸",
		"IT": "🇮🇹",
		"JE": "🇯🇪",
		"JM": "🇯🇲",
		"JO": "🇯🇴",
		"JP": "🇯🇵",
		"KE": "🇰🇪",
		"KG": "🇰🇬",
		"KH": "🇰🇭",
		"KI": "🇰🇮",
		"KM": "🇰🇲",
		"KN": "🇰🇳",
		"KP": "🇰🇵",
		"KR": "🇰🇷",
		"KW": "🇰🇼",
		"KY": "🇰🇾",
		"KZ": "🇰🇿",
		"LA": "🇱🇦",
		"LB": "🇱🇧",
		"LC": "🇱🇨",
		"LI": "🇱🇮",
		"LK": "🇱🇰",
		"LR": "🇱🇷",
		"LS": "🇱🇸",
		"LT": "🇱🇹",
		"LU": "🇱🇺",
		"LV": "🇱🇻",
		"LY": "🇱🇾",
		"MA": "🇲🇦",
		"MC": "🇲🇨",
		"MD": "🇲🇩",
		"ME": "🇲🇪",
		"MF": "🇲🇫",
		"MG": "🇲🇬",
		"MH": "🇲🇭",
		"MK": "🇲🇰",
		"ML": "🇲🇱",
		"MM": "🇲🇲",
		"MN": "🇲🇳",
		"MO": "🇲🇴",
		"MP": "🇲🇵",
		"MQ": "🇲🇶",
		"MR": "🇲🇷",
		"MS": "🇲🇸",
		"MT": "🇲🇹",
		"MU": "🇲🇺",
		"MV": "🇲🇻",
		"MW": "🇲🇼",
		"MX": "🇲🇽",
		"MY": "🇲🇾",
		"MZ": "🇲🇿",
		"NA": "🇳🇦",
		"NC": "🇳🇨",
		"NE": "🇳🇪",
		"NF": "🇳🇫",
		"NG": "🇳🇬",
		"NI": "🇳🇮",
		"NL": "🇳🇱",
		"NO": "🇳🇴",
		"NP": "🇳🇵",
		"NR": "🇳🇷",
		"NU": "🇳🇺",
		"NZ": "🇳🇿",
		"OM": "🇴🇲",
		"PA": "🇵🇦",
		"PE": "🇵🇪",
		"PF": "🇵🇫",
		"PG": "🇵🇬",
		"PH": "🇵🇭",
		"PK": "🇵🇰",
		"PL": "🇵🇱",
		"PM": "🇵🇲",
		"PN": "🇵🇳",
		"PR": "🇵🇷",
		"PS": "🇵🇸",
		"PT": "🇵🇹",
		"PW": "🇵🇼",
		"PY": "🇵🇾",
		"QA": "🇶🇦",
		"RE": "🇷🇪",
		"RO": "🇷🇴",
		"RS": "🇷🇸",
		"RU": "🇷🇺",
		"RW": "🇷🇼",
		"SA": "🇸🇦",
		"SB": "🇸🇧",
		"SC": "🇸🇨",
		"SD": "🇸🇩",
		"SE": "🇸🇪",
		"SG": "🇸🇬",
		"SH": "🇸🇭",
		"SI": "🇸🇮",
		"SJ": "🇸🇯",
		"SK": "🇸🇰",
		"SL": "🇸🇱",
		"SM": "🇸🇲",
		"SN": "🇸🇳",
		"SO": "🇸🇴",
		"SR": "🇸🇷",
		"SS": "🇸🇸",
		"ST": "🇸🇹",
		"SV": "🇸🇻",
		"SX": "🇸🇽",
		"SY": "🇸🇾",
		"SZ": "🇸🇿",
		"TC": "🇹🇨",
		"TD": "🇹🇩",
		"TF": "🇹🇫",
		"TG": "🇹🇬",
		"TH": "🇹🇭",
		"TJ": "🇹🇯",
		"TK": "🇹🇰",
		"TL": "🇹🇱",
		"TM": "🇹🇲",
		"TN": "🇹🇳",
		"TO": "🇹🇴",
		"TR": "🇹🇷",
		"TT": "🇹🇹",
		"TV": "🇹🇻",
		"TW": "🇹🇼",
		"TZ": "🇹🇿",
		"UA": "🇺🇦",
		"UG": "🇺🇬",
		"UM": "🇺🇲",
		"US": "🇺🇸",
		"UY": "🇺🇾",
		"UZ": "🇺🇿",
		"VA": "🇻🇦",
		"VC": "🇻🇨",
		"VE": "🇻🇪",
		"VG": "🇻🇬",
		"VI": "🇻🇮",
		"VN": "🇻🇳",
		"VU": "🇻🇺",
		"WF": "🇼🇫",
		"WS": "🇼🇸",
		"YE": "🇾🇪",
		"YT": "🇾🇹",
		"ZA": "🇿🇦",
		"ZM": "🇿🇲",
		"ZW": "🇿🇼",
	}

	flag, ok := flags[code]
	if !ok {
		return ""
	}
	return flag
}
