package main

import (
	"github.com/mlange-42/ark-pixel/plot"
	"github.com/mlange-42/ark-pixel/window"
	"github.com/mlange-42/ark-tools/app"
)

func main() {
	app := app.New()
	app.TPS = 30
	app.AddUISystem((&window.Window{}).
		With(&plot.PerfStats{}).
		With(&plot.Controls{}))

	app.AddUISystem((&window.Window{}).
		With(&plot.Systems{}))

	//app.AddSystem(&system.FixedTermination{
	//	Steps: 100,
	//})

	window.Run(app)
}
