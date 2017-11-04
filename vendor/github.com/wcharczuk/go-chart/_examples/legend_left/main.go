package main

import (
	"net/http"

	"github.com/wcharczuk/go-chart"
)

func drawChart(res http.ResponseWriter, req *http.Request) {

	/*
	   In this example we add a `Renderable` or a custom component to the `Elements` array.
	   In this specific case it is a pre-built renderable (`CreateLegend`) that draws a legend for the chart's series.
	   If you like, you can use `CreateLegend` as a template for writing your own renderable, or even your own legend.
	*/

	graph := chart.Chart{
		XAxis: chart.XAxis{
			Style: chart.Style{Show: true},
		},
		YAxis: chart.YAxis{
			Style: chart.Style{Show: true},
		},
		Background: chart.Style{
			Padding: chart.Box{
				Top:  20,
				Left: 260,
			},
		},
		Series: []chart.Series{
			chart.ContinuousSeries{
				Name:    "A test series",
				XValues: []float64{1.0, 2.0, 3.0, 4.0, 5.0},
				YValues: []float64{1.0, 2.0, 3.0, 4.0, 5.0},
			},

			chart.ContinuousSeries{
				Name:    "Another test series",
				XValues: []float64{1.0, 2.0, 3.0, 4.0, 5.0},
				YValues: []float64{1.0, 2.0, 3.0, 4.0, 5.0},
			},

			chart.ContinuousSeries{
				Name:    "Yet Another test series",
				XValues: []float64{1.0, 2.0, 3.0, 4.0, 5.0},
				YValues: []float64{1.0, 2.0, 3.0, 4.0, 5.0},
			},

			chart.ContinuousSeries{
				Name:    "Even More series",
				XValues: []float64{1.0, 2.0, 3.0, 4.0, 5.0},
				YValues: []float64{1.0, 2.0, 3.0, 4.0, 5.0},
			},

			chart.ContinuousSeries{
				Name:    "Foo Bar",
				XValues: []float64{1.0, 2.0, 3.0, 4.0, 5.0},
				YValues: []float64{1.0, 2.0, 3.0, 4.0, 5.0},
			},

			chart.ContinuousSeries{
				Name:    "Bar Baz",
				XValues: []float64{1.0, 2.0, 3.0, 4.0, 5.0},
				YValues: []float64{1.0, 2.0, 3.0, 4.0, 5.0},
			},

			chart.ContinuousSeries{
				Name:    "Moo Bar",
				XValues: []float64{1.0, 2.0, 3.0, 4.0, 5.0},
				YValues: []float64{1.0, 2.0, 3.0, 4.0, 5.0},
			},

			chart.ContinuousSeries{
				Name:    "Zoo Bar Baz",
				XValues: []float64{1.0, 2.0, 3.0, 4.0, 5.0},
				YValues: []float64{1.0, 2.0, 3.0, 4.0, 5.0},
			},

			chart.ContinuousSeries{
				Name:    "Fast and the Furious",
				XValues: []float64{1.0, 2.0, 3.0, 4.0, 5.0},
				YValues: []float64{5.0, 4.0, 3.0, 2.0, 1.0},
			},

			chart.ContinuousSeries{
				Name:    "2 Fast 2 Furious",
				XValues: []float64{1.0, 2.0, 3.0, 4.0, 5.0},
				YValues: []float64{5.0, 4.0, 3.0, 2.0, 1.0},
			},

			chart.ContinuousSeries{
				Name:    "They only get more fast and more furious",
				XValues: []float64{1.0, 2.0, 3.0, 4.0, 5.0},
				YValues: []float64{5.0, 4.0, 3.0, 2.0, 1.0},
			},
		},
	}

	//note we have to do this as a separate step because we need a reference to graph
	graph.Elements = []chart.Renderable{
		chart.LegendLeft(&graph),
	}

	res.Header().Set("Content-Type", "image/png")
	graph.Render(chart.PNG, res)
}

func main() {
	http.HandleFunc("/", drawChart)
	http.ListenAndServe(":8080", nil)
}
