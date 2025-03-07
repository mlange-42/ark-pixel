package main

import (
	"github.com/mlange-42/ark-pixel/monitor"
	"github.com/mlange-42/ark-pixel/window"
	"github.com/mlange-42/ark-tools/app"
)

func main() {
	app := app.New()
	app.TPS = 30

	app.AddUISystem((&window.Window{}).
		With(&monitor.PerfStats{}).
		With(&monitor.Controls{}))

	app.AddUISystem((&window.Window{}).
		With(&monitor.Systems{}))

	app.AddUISystem((&window.Window{}).
		With(&monitor.Resources{}))

	window.Run(app)
}
