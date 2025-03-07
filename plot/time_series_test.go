package plot_test

import (
	"math/rand"
	"testing"

	"github.com/mlange-42/ark-pixel/plot"
	"github.com/mlange-42/ark-pixel/window"
	"github.com/mlange-42/ark-tools/app"
	"github.com/mlange-42/ark-tools/system"
	"github.com/mlange-42/ark/ecs"
)

func ExampleTimeSeries() {

	// Create a new model.
	app := app.New()

	// Limit the the simulation speed.
	app.TPS = 30

	// Create a time series plot.
	// See below for the implementation of the RowObserver.
	app.AddUISystem((&window.Window{}).
		With(&plot.TimeSeries{
			Observer: &RowObserver{},
		}))

	// Add a termination system that ends the simulation.
	app.AddSystem(&system.FixedTermination{
		Steps: 100,
	})

	// Run the simulation.
	// Due to the use of the OpenGL UI system, the model must be run via [window.Run].
	// Uncomment the next line to run this example stand-alone.

	window.Run(app)

	// Output:
}

func TestTimeSeries_Columns(t *testing.T) {
	app := app.New()
	app.TPS = 300
	app.AddUISystem((&window.Window{}).
		With(&plot.TimeSeries{
			Observer: &RowObserver{},
			Columns:  []string{"A", "C"},
			MaxRows:  50,
		}))

	app.AddSystem(&system.FixedTermination{
		Steps: 100,
	})
	app.Run()
}

func TestTimeSeries_PanicColumns(t *testing.T) {
	app := app.New()
	app.TPS = 300
	app.AddUISystem((&window.Window{}).
		With(&plot.TimeSeries{
			Observer: &RowObserver{},
			Columns:  []string{"A", "F"},
		}))

	app.AddSystem(&system.FixedTermination{
		Steps: 100,
	})
	//assert.Panics(t, app.Run)
}

// RowObserver to generate random time series.
type RowObserver struct{}

func (o *RowObserver) Initialize(w *ecs.World) {}
func (o *RowObserver) Update(w *ecs.World)     {}
func (o *RowObserver) Header() []string {
	return []string{"A", "B", "C"}
}
func (o *RowObserver) Values(w *ecs.World) []float64 {
	return []float64{rand.Float64(), rand.Float64() + 1, rand.Float64() + 2}
}
