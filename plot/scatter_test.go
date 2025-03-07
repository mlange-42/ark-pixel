package plot_test

import (
	"testing"

	"github.com/mlange-42/ark-pixel/plot"
	"github.com/mlange-42/ark-pixel/window"
	"github.com/mlange-42/ark-tools/app"
	"github.com/mlange-42/ark-tools/observer"
	"github.com/mlange-42/ark-tools/system"
	"github.com/stretchr/testify/assert"
)

func ExampleScatter() {

	// Create a new model.
	app := app.New()

	// Limit the the simulation speed.
	app.TPS = 30

	// Create a scatter plot.
	app.AddUISystem((&window.Window{}).
		With(&plot.Scatter{
			Observers: []observer.Table{
				&TableObserver{}, // One or more observers.
			},
			X: []string{
				"X", // One X column per observer.
			},
			Y: [][]string{
				{"A", "B", "C"}, // One or more Y columns per observer.
			},
		}))

	// Add a termination system that ends the simulation.
	app.AddSystem(&system.FixedTermination{
		Steps: 100,
	})

	// Run the simulation.
	// Due to the use of the OpenGL UI system, the model must be run via [window.Run].
	// Uncomment the next line to run this example stand-alone.

	//window.Run(app)

	// Output:
}

func TestScatter(t *testing.T) {
	app := app.New()
	app.TPS = 300

	app.AddUISystem((&window.Window{}).
		With(&plot.Scatter{
			Observers: []observer.Table{
				&TableObserver{},
			},
		}))

	app.AddSystem(&system.FixedTermination{
		Steps: 100,
	})
	//app.Run()
}

func TestScatter_PanicXCount(t *testing.T) {
	app := app.New()
	app.TPS = 300

	app.AddUISystem((&window.Window{}).
		With(&plot.Scatter{
			Observers: []observer.Table{
				&TableObserver{},
			},
			X: []string{"X", "X"},
		}))

	app.AddSystem(&system.FixedTermination{
		Steps: 100,
	})
	assert.Panics(t, app.Run)
}

func TestScatter_PanicX(t *testing.T) {
	app := app.New()
	app.TPS = 300

	app.AddUISystem((&window.Window{}).
		With(&plot.Scatter{
			Observers: []observer.Table{
				&TableObserver{},
			},
			X: []string{"F"},
		}))

	app.AddSystem(&system.FixedTermination{
		Steps: 100,
	})
	assert.Panics(t, app.Run)
}

func TestScatter_PanicYCount(t *testing.T) {
	app := app.New()
	app.TPS = 300

	app.AddUISystem((&window.Window{}).
		With(&plot.Scatter{
			Observers: []observer.Table{
				&TableObserver{},
			},
			Y: [][]string{
				{"A"},
				{"A"},
			},
		}))

	app.AddSystem(&system.FixedTermination{
		Steps: 100,
	})
	assert.Panics(t, app.Run)
}

func TestScatter_PanicY(t *testing.T) {
	app := app.New()
	app.TPS = 300

	app.AddUISystem((&window.Window{}).
		With(&plot.Scatter{
			Observers: []observer.Table{
				&TableObserver{},
			},
			Y: [][]string{
				{"F"},
			},
		}))

	app.AddSystem(&system.FixedTermination{
		Steps: 100,
	})
	assert.Panics(t, app.Run)
}
