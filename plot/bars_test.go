package plot_test

import (
	"testing"

	"github.com/mlange-42/ark-pixel/plot"
	"github.com/mlange-42/ark-pixel/window"
	"github.com/mlange-42/ark-tools/app"
	"github.com/mlange-42/ark-tools/system"
	"github.com/stretchr/testify/assert"
)

func ExampleBars() {

	// Create a new model.
	app := app.New()

	// Limit the the simulation speed.
	app.TPS = 30

	// Create a time series plot.
	app.AddUISystem((&window.Window{}).
		With(&plot.Bars{
			Observer: &RowObserver{},
			YLim:     [...]float64{0, 4}, // Optional Y axis limits.
		}))

	// Add a termination system that ends the simulation.
	app.AddSystem(&system.FixedTermination{
		Steps: 100,
	})

	app.Run()

	// Run the simulation.
	// Due to the use of the OpenGL UI system, the model must be run via [window.Run].
	// Comment out the code line above, and uncomment the next line to run this example stand-alone.

	// window.Run(m)

	// Output:
}

func TestBars_Columns(t *testing.T) {
	app := app.New()
	app.TPS = 300
	app.AddUISystem((&window.Window{}).
		With(&plot.Bars{
			Observer: &RowObserver{},
			YLim:     [...]float64{0, 4},
			Columns:  []string{"A", "C"},
		}))
	app.AddSystem(&system.FixedTermination{
		Steps: 100,
	})
	//app.Run()
}

func TestBars_PanicColumns(t *testing.T) {
	app := app.New()
	app.TPS = 300
	app.AddUISystem((&window.Window{}).
		With(&plot.Bars{
			Observer: &RowObserver{},
			Columns:  []string{"A", "F"},
		}))
	app.AddSystem(&system.FixedTermination{
		Steps: 100,
	})
	assert.Panics(t, app.Run)
}
