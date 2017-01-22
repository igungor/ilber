package main

import (
	"net/http"

	"github.com/wcharczuk/go-chart"
)

func drawChart(res http.ResponseWriter, req *http.Request) {
	mainSeries := chart.ContinuousSeries{
		Name:    "A test series",
		XValues: chart.Sequence.Float64(1.0, 100.0),
		YValues: chart.Sequence.RandomWithAverage(100, 100, 50),
	}

	minSeries := &chart.MinSeries{
		Style: chart.Style{
			Show:            true,
			StrokeColor:     chart.ColorAlternateGray,
			StrokeDashArray: []float64{5.0, 5.0},
		},
		InnerSeries: mainSeries,
	}

	maxSeries := &chart.MaxSeries{
		Style: chart.Style{
			Show:            true,
			StrokeColor:     chart.ColorAlternateGray,
			StrokeDashArray: []float64{5.0, 5.0},
		},
		InnerSeries: mainSeries,
	}

	graph := chart.Chart{
		Width:  1920,
		Height: 1080,
		YAxis: chart.YAxis{
			Name:      "Random Values",
			NameStyle: chart.StyleShow(),
			Style:     chart.StyleShow(),
			Range: &chart.ContinuousRange{
				Min: 25,
				Max: 175,
			},
		},
		XAxis: chart.XAxis{
			Name:      "Random Other Values",
			NameStyle: chart.StyleShow(),
			Style:     chart.StyleShow(),
		},
		Series: []chart.Series{
			mainSeries,
			minSeries,
			maxSeries,
			chart.LastValueAnnotation(minSeries),
			chart.LastValueAnnotation(maxSeries),
		},
	}

	graph.Elements = []chart.Renderable{chart.Legend(&graph)}

	res.Header().Set("Content-Type", "image/png")
	graph.Render(chart.PNG, res)
}

func main() {
	http.HandleFunc("/", drawChart)
	http.ListenAndServe(":8080", nil)
}
