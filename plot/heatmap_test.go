package plot_test

import (
	"github.com/mlange-42/ark-pixel/plot"
	"github.com/mlange-42/ark-pixel/window"
	"github.com/mlange-42/ark-tools/app"
	"github.com/mlange-42/ark-tools/observer"
	"github.com/mlange-42/ark-tools/system"
	"gonum.org/v1/plot/palette"
)

func ExampleHeatMap() {
	// Create a new model.
	app := app.New()

	// Limit the the simulation speed.
	app.TPS = 30
	app.FPS = 0

	// Create a contour plot.
	app.AddUISystem(
		(&window.Window{}).
			With(&plot.HeatMap{
				Observer: observer.MatrixToGrid(&MatrixObserver{}, nil, nil),
				Palette:  palette.Heat(16, 1),
				Min:      -2,
				Max:      2,
			}))

	// Add a termination system that ends the simulation.
	app.AddSystem(&system.FixedTermination{
		Steps: 100,
	})

	// Run the simulation.
	// Due to the use of the OpenGL UI system, the model must be run via [window.Run].
	// Uncomment the next line to run this example stand-alone.

	// window.Run(m)

	// Output:
}
