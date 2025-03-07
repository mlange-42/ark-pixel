package plot_test

import (
	"math"
	"testing"

	"github.com/mazznoer/colorgrad"
	"github.com/mlange-42/ark-pixel/plot"
	"github.com/mlange-42/ark-pixel/window"
	"github.com/mlange-42/ark-tools/app"
	"github.com/mlange-42/ark-tools/system"
	"github.com/mlange-42/ark/ecs"
)

func ExampleImage() {

	// Create a new model.
	app := app.New()

	// Limit the the simulation speed.
	app.TPS = 30
	app.FPS = 0

	// Create an image plot.
	// See below for the implementation of the MatrixObserver.
	app.AddUISystem(
		(&window.Window{}).
			With(&plot.Image{
				Scale:    4,
				Observer: &MatrixObserver{},
				Colors:   colorgrad.Inferno(),
				Min:      -2,
				Max:      2,
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

func TestImage_LimitsScale(t *testing.T) {
	app := app.New()
	app.TPS = 300
	app.FPS = 0
	app.AddUISystem(
		(&window.Window{}).
			With(&plot.Image{
				Observer: &MatrixObserver{},
				Colors:   colorgrad.Inferno(),
			}))

	app.AddSystem(&system.FixedTermination{
		Steps: 100,
	})

	//app.Run()
}

// Example observer, reporting a matrix with z = sin(0.1*i) + sin(0.2*j).
type MatrixObserver struct {
	cols   int
	rows   int
	values []float64
}

func (o *MatrixObserver) Initialize(w *ecs.World) {
	o.cols = 160
	o.rows = 120
	o.values = make([]float64, o.cols*o.rows)
}

func (o *MatrixObserver) Update(w *ecs.World) {}

func (o *MatrixObserver) Dims() (int, int) {
	return o.cols, o.rows
}

func (o *MatrixObserver) Values(w *ecs.World) []float64 {
	for idx := 0; idx < len(o.values); idx++ {
		i := idx % o.cols
		j := idx / o.cols
		o.values[idx] = math.Sin(0.1*float64(i)) + math.Sin(0.2*float64(j))
	}
	return o.values
}
