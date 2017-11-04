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

const (
	defaultDataRange = "3d"
)

func init() {
	register(cmdChart)
}

var cmdChart = &Command{
	Name:      "chart",
	ShortLine: "chart kaba kaat",
	Run:       runChart,
	Hidden:    true,
}

func runChart(ctx context.Context, b *bot.Bot, msg *telegram.Message) {
	args := msg.Args()

	var currencies []string
	switch len(args) {
	case 0:
		currencies = []string{"USD", "TRY"}
	case 1:
		currencies = []string{normalize(args[0]), "TRY"}
	case 2:
		currencies = []string{normalize(args[0]), normalize(args[1])}
	default:
		_, _ = b.SendMessage(msg.Chat.ID, "anlamadim")
		return
	}

	u, _ := url.Parse(yahooFinanceURL)
	u.Path += fmt.Sprintf("%v%v=%v", currencies[0], currencies[1], "X")
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
		dur, _ := time.ParseDuration(defaultDataRange)
		if dur < 2*24*time.Hour {
			return chart.TimeHourValueFormatter(v)
		}
		return chart.TimeValueFormatter(v)
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

	chartname := fmt.Sprintf("%v in %v for %v", currencies[0], currencies[1], defaultDataRange)
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
