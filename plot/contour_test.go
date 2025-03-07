package plot_test

import (
	"testing"

	"github.com/mlange-42/ark-pixel/plot"
	"github.com/mlange-42/ark-pixel/window"
	"github.com/mlange-42/ark-tools/app"
	"github.com/mlange-42/ark-tools/observer"
	"github.com/mlange-42/ark-tools/system"
	"gonum.org/v1/plot/palette"
)

func ExampleContour() {
	// Create a new model.
	app := app.New()

	// Limit the the simulation speed.
	app.TPS = 30
	app.FPS = 0

	// Create a contour plot.
	app.AddUISystem(
		(&window.Window{}).
			With(&plot.Contour{
				Observer: observer.MatrixToGrid(&MatrixObserver{}, nil, nil),
				Palette:  palette.Heat(16, 1),
				Levels:   []float64{-2, -1.5, -1, -0.5, 0, 0.5, 1, 1.5, 2},
			}))

	// Add a termination system that ends the simulation.
	app.AddSystem(&system.FixedTermination{
		Steps: 100,
	})

	app.Run()

	// Run the simulation.
	// Due to the use of the OpenGL UI system, the model must be run via [window.Run].
	// Comment out the code line above, and uncomment the next line to run this example stand-alone.

	// window.Run(app)

	// Output:
}

func TestContour_NoLevels(t *testing.T) {
	app := app.New()
	app.TPS = 300
	app.FPS = 0

	app.AddUISystem(
		(&window.Window{}).
			With(&plot.Contour{
				Observer: observer.MatrixToGrid(&MatrixObserver{}, nil, nil),
				Palette:  palette.Heat(16, 1),
			}))

	app.AddSystem(&system.FixedTermination{
		Steps: 100,
	})

	app.Run()
}
