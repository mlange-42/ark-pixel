package plot_test

import (
	"testing"

	"github.com/mlange-42/ark-pixel/plot"
	"github.com/mlange-42/ark-pixel/window"
	"github.com/mlange-42/ark-tools/app"
	"github.com/mlange-42/ark-tools/system"
	"github.com/stretchr/testify/assert"
)

func ExampleLines() {

	// Create a new model.
	app := app.New()

	// Limit the the simulation speed.
	app.TPS = 30

	// Create a line plot.
	// See below for the implementation of the TableObserver.
	app.AddUISystem((&window.Window{}).
		With(&plot.Lines{
			Observer: &TableObserver{},
			X:        "X",                     // Optional, defaults to row index
			Y:        []string{"A", "B", "C"}, // Optional, defaults to all but X
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

func TestLines(t *testing.T) {
	app := app.New()
	app.TPS = 300
	app.AddUISystem((&window.Window{}).
		With(&plot.Lines{
			Observer: &TableObserver{},
			XLim:     [2]float64{0, 30},
			YLim:     [2]float64{0.5, 0.6},
		}))

	app.AddSystem(&system.FixedTermination{
		Steps: 100,
	})
	app.Run()
}

func TestLines_PanicX(t *testing.T) {
	app := app.New()
	app.AddUISystem((&window.Window{}).
		With(&plot.Lines{
			Observer: &TableObserver{},
			X:        "U",
		}))

	app.AddSystem(&system.FixedTermination{
		Steps: 100,
	})
	assert.Panics(t, app.Run)
}

func TestLines_PanicY(t *testing.T) {
	app := app.New()
	app.AddUISystem((&window.Window{}).
		With(&plot.Lines{
			Observer: &TableObserver{},
			X:        "X",
			Y:        []string{"A", "B", "U"},
		}))

	app.AddSystem(&system.FixedTermination{
		Steps: 100,
	})
	assert.Panics(t, app.Run)
}

func TestLinesNaN(t *testing.T) {
	app := app.New()
	app.AddUISystem((&window.Window{}).
		With(&plot.Lines{
			Observer: &TableObserverNaN{},
			X:        "X",
			Y:        []string{"A"},
		}))

	app.AddSystem(&system.FixedTermination{
		Steps: 100,
	})
	app.Run()
}
