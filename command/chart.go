package command

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/igungor/ilber/bot"
	"github.com/igungor/telegram"
	"github.com/wcharczuk/go-chart"
	"github.com/wcharczuk/go-chart/drawing"
)

func init() {
	// 	register(cmdChart)
}

var cmdChart = &Command{
	Name:      "chart",
	ShortLine: "bişeyler çiziktir",
	Run:       runChart,
	Hidden:    true,
}

func runChart(ctx context.Context, b *bot.Bot, msg *telegram.Message) {
	args := msg.Args()
	var opts telegram.SendOptions

	currency := "dollar"
	if len(args) > 0 {
		currency = args[0]
	}

	pairs, err := b.Store.Values(currency)
	if err != nil {
		_, _ = b.SendMessage(msg.Chat.ID, err.Error(), &opts)
		return
	}

	mainSeries := chart.ContinuousSeries{
		Name:    "currency",
		XValues: chart.Sequence.Float64(1, 100),
		YValues: chart.Sequence.RandomWithAverage(100, 100, 50),
	}

	maxSeries := &chart.MaxSeries{
		Name: "max value",
		Style: chart.Style{
			Show:            true,
			StrokeColor:     chart.ColorAlternateGray,
			StrokeDashArray: []float64{5, 5},
		},
		InnerSeries: mainSeries,
	}

	avgSeries := chart.SMASeries{
		Name: "Average",
		Style: chart.Style{
			Show:            true,
			StrokeColor:     drawing.ColorRed,
			StrokeDashArray: []float64{5, 5},
		},
		InnerSeries: mainSeries,
	}

	graph := chart.Chart{
		XAxis: chart.XAxis{
			Name:      "Random Other Values",
			NameStyle: chart.StyleShow(),
			Style:     chart.StyleShow(),
		},
		YAxis: chart.YAxis{
			Name:      "Random Values",
			NameStyle: chart.StyleShow(),
			Style:     chart.StyleShow(),
			Range: &chart.ContinuousRange{
				Min: 25,
				Max: 175,
			},
		},
		Series: []chart.Series{
			mainSeries,
			maxSeries,
			avgSeries,
			chart.LastValueAnnotation(maxSeries),
		},
	}

	graph.Elements = []chart.Renderable{chart.Legend(&graph)}

	var buf bytes.Buffer
	err = graph.Render(chart.PNG, &buf)
	if err != nil {
		log.Printf("Error rendering image: %v\n", err)
		return
	}

	chartname := fmt.Sprintf("%v-%v.png", currency, time.Now().Format("2006-01-02-15:04:05"))
	photo := telegram.Photo{File: telegram.File{
		Name: chartname,
		Body: &buf,
	}}
	_, err = b.SendPhoto(msg.Chat.ID, photo, nil)
	if err != nil {
		log.Printf("Error sending photo: %v\n", err)
		return
	}
}
