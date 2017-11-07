package command

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/igungor/ilber/bot"
	"github.com/igungor/telegram"
	chart "github.com/wcharczuk/go-chart"
)

func init() {
	register(cmdChart)
}

const (
	defaultDataRange = "5d"
)

var cmdChart = &Command{
	Name:      "chart",
	ShortLine: "chart kaba kaat",
	Run:       runChart,
	Hidden:    true,
}

func runChart(ctx context.Context, b *bot.Bot, msg *telegram.Message) {
	args := msg.Args()

	var (
		fromCurrency string
		toCurrency   string
	)
	switch len(args) {
	case 0:
		fromCurrency = "USD"
		toCurrency = "TRY"
	case 1:
		fromCurrency = normalize(args[0])
		toCurrency = "TRY"
	case 2:
		fromCurrency = normalize(args[0])
		toCurrency = normalize(args[1])
	default:
		_, _ = b.SendMessage(msg.Chat.ID, "anlamadim")
		return
	}

	u, _ := url.Parse(yahooFinanceURL)
	u.Path += fmt.Sprintf("%v%v=%v", fromCurrency, toCurrency, "X")
	params := u.Query()
	params.Set("range", defaultDataRange)
	u.RawQuery = params.Encode()

	resp, err := httpclient.Get(u.String())
	if err != nil {
		b.Logger.Printf("chart: could not fetch response: %v", err)
		b.SendMessage(msg.Chat.ID, "bir takim hatalar sozkonusu")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b.Logger.Printf("chart: status not OK: %v", resp.StatusCode)
		b.SendMessage(msg.Chat.ID, "bir takim hatalar sozkonusu")
		return
	}

	var response yahooFinanceResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		b.Logger.Printf("chart: could not parse json: %v", err)
		b.SendMessage(msg.Chat.ID, "bir takim hatalar sozkonusu")
		return
	}

	result := response.Chart.Result
	quote := result[len(result)-1].Indicators.Quote
	if len(quote) == 0 {
		b.Logger.Printf("chart: no value found for query %v", strings.Join(args, " "))
		b.SendMessage(msg.Chat.ID, "no value found")
		return
	}
	close := quote[len(quote)-1].Close
	close = close[:len(close)-1] // last value is always nil
	timestamps := result[0].Timestamp
	timestamps = timestamps[:len(timestamps)-1] // last value of the rate is nil, so remove the last timestamp as well
	var rates []float64
	ignoredidx := make([]bool, len(close))
	for i, v := range close {
		rate, ok := v.(float64)
		// skip unrecognized values to a list for later use
		if !ok {
			ignoredidx[i] = true
			continue
		}
		rates = append(rates, rate)
	}

	var times []time.Time
	for i, ts := range timestamps {
		// if the rate value is nil for this timestamp, skip it so that the
		// data will be consistent
		if ignoredidx[i] {
			continue
		}
		times = append(times, time.Unix(ts, 0))
	}

	floatformatter := func(v interface{}) string {
		return fmt.Sprintf("%4.4f", v)
	}

	dateformatter := func(v interface{}) string {
		duration := toDuration(defaultDataRange)
		if duration <= 2*24*time.Hour {
			return chart.TimeHourValueFormatter(v)
		}
		return chart.TimeValueFormatterWithFormat("01-02 03:04PM")(v)
	}

	priceSeries := chart.TimeSeries{
		Name: "Currency",
		Style: chart.Style{
			Show:        true,
			StrokeColor: chart.ColorBlue,
		},
		XValues: times,
		YValues: rates,
	}

	chartname := fmt.Sprintf("%v in %v for %v", fromCurrency, toCurrency, defaultDataRange)
	graph := chart.Chart{
		Title:      chartname,
		TitleStyle: chart.StyleShow(),
		XAxis: chart.XAxis{
			Style:          chart.StyleShow(),
			ValueFormatter: dateformatter,
			TickStyle:      chart.StyleShow(),
			GridMajorStyle: chart.Style{
				Show:            true,
				StrokeColor:     chart.ColorAlternateGray,
				StrokeDashArray: []float64{3, 2, 1},
				StrokeWidth:     1.0,
			},
		},
		YAxis: chart.YAxis{
			Name:           "Currency",
			NameStyle:      chart.StyleShow(),
			Style:          chart.StyleShow(),
			ValueFormatter: floatformatter,
		},
		Series: []chart.Series{
			priceSeries,
			chart.LastValueAnnotation(priceSeries, floatformatter),
		},
	}

	var buf bytes.Buffer
	graph.Elements = []chart.Renderable{chart.Legend(&graph)}
	if err := graph.Render(chart.PNG, &buf); err != nil {
		b.Logger.Printf("chart: could not render the chart: %v", err)
		_, _ = b.SendMessage(msg.Chat.ID, "hersey tamamdi aslinda ama grafigi render edemedim")
		return
	}

	photo := telegram.Photo{File: telegram.File{
		Name: chartname,
		Body: &buf,
	}}

	_, err = b.SendPhoto(msg.Chat.ID, photo, nil)
	if err != nil {
		b.Logger.Printf("chart: could not send photo: %v", err)
		return
	}
}

func toDuration(s string) time.Duration {
	var (
		day   = 24 * time.Hour
		month = 30 * day
		year  = 12 * month

		defaultDuration = 5 * day
	)
	validRanges := map[string]time.Duration{
		"1d":  day,
		"5d":  5 * day,
		"1mo": month,
		"3mo": 3 * month,
		"6mo": 6 * month,
		"1y":  year,
		"2y":  2 * year,
		"5y":  5 * year,
		"10y": 10 * year,
	}
	if dur, ok := validRanges[s]; ok {
		return dur
	}
	return defaultDuration
}
