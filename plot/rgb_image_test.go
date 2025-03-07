package plot_test

import (
	"math"
	"testing"

	"github.com/mlange-42/ark-pixel/plot"
	"github.com/mlange-42/ark-pixel/window"
	"github.com/mlange-42/ark-tools/app"
	"github.com/mlange-42/ark-tools/observer"
	"github.com/mlange-42/ark-tools/system"
	"github.com/mlange-42/ark/ecs"
	"github.com/stretchr/testify/assert"
)

func ExampleImageRGB() {

	// Create a new model.
	app := app.New()

	// Limit the the simulation speed.
	app.TPS = 30

	// Create an RGB image plot.
	// See below for the implementation of the CallbackMatrixObserver.
	app.AddUISystem((&window.Window{}).
		With(&plot.ImageRGB{
			Observer: observer.MatrixToLayers(
				&CallbackMatrixObserver{Callback: func(i, j int) float64 { return float64(i) / 240 }},
				&CallbackMatrixObserver{Callback: func(i, j int) float64 { return math.Sin(0.1 * float64(i)) }},
				&CallbackMatrixObserver{Callback: func(i, j int) float64 { return float64(j) / 160 }},
			),
			Min: []float64{0, 0, 0},
			Max: []float64{1, 1, 1},
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

func TestImageRGB(t *testing.T) {
	app := app.New()
	app.TPS = 300
	app.AddUISystem((&window.Window{}).
		With(&plot.ImageRGB{
			Observer: observer.MatrixToLayers(
				&CallbackMatrixObserver{Callback: func(i, j int) float64 { return float64(i) / 240 }},
				&CallbackMatrixObserver{Callback: func(i, j int) float64 { return math.Sin(0.1 * float64(i)) }},
				&CallbackMatrixObserver{Callback: func(i, j int) float64 { return float64(j) / 160 }},
			),
		}))
	app.AddSystem(&system.FixedTermination{
		Steps: 100,
	})
	app.Run()
}

func TestImageRGB_PanicMin(t *testing.T) {
	app := app.New()
	app.TPS = 300
	app.AddUISystem((&window.Window{}).
		With(&plot.ImageRGB{
			Observer: observer.MatrixToLayers(
				&CallbackMatrixObserver{Callback: func(i, j int) float64 { return float64(i) / 240 }},
				&CallbackMatrixObserver{Callback: func(i, j int) float64 { return math.Sin(0.1 * float64(i)) }},
				&CallbackMatrixObserver{Callback: func(i, j int) float64 { return float64(j) / 160 }},
			),
			Min: []float64{0, 0},
		}))
	app.AddSystem(&system.FixedTermination{
		Steps: 100,
	})
	assert.Panics(t, app.Run)
}

func TestImageRGB_PanicMax(t *testing.T) {
	app := app.New()
	app.TPS = 300
	app.AddUISystem((&window.Window{}).
		With(&plot.ImageRGB{
			Observer: observer.MatrixToLayers(
				&CallbackMatrixObserver{Callback: func(i, j int) float64 { return float64(i) / 240 }},
				&CallbackMatrixObserver{Callback: func(i, j int) float64 { return math.Sin(0.1 * float64(i)) }},
				&CallbackMatrixObserver{Callback: func(i, j int) float64 { return float64(j) / 160 }},
			),
			Max: []float64{1, 1},
		}))
	app.AddSystem(&system.FixedTermination{
		Steps: 100,
	})
	assert.Panics(t, app.Run)
}

func TestImageRGB_PanicLayerCount(t *testing.T) {
	app := app.New()
	app.TPS = 300
	app.AddUISystem((&window.Window{}).
		With(&plot.ImageRGB{
			Observer: observer.MatrixToLayers(
				&CallbackMatrixObserver{Callback: func(i, j int) float64 { return float64(i) / 240 }},
				&CallbackMatrixObserver{Callback: func(i, j int) float64 { return math.Sin(0.1 * float64(i)) }},
				&CallbackMatrixObserver{Callback: func(i, j int) float64 { return float64(j) / 160 }},
			),
			Layers: []int{2, 1, 2, 0},
		}))
	app.AddSystem(&system.FixedTermination{
		Steps: 100,
	})
	assert.Panics(t, app.Run)
}

func TestImageRGB_PanicLayerIndex(t *testing.T) {
	app := app.New()
	app.TPS = 300
	app.AddUISystem((&window.Window{}).
		With(&plot.ImageRGB{
			Observer: observer.MatrixToLayers(
				&CallbackMatrixObserver{Callback: func(i, j int) float64 { return float64(i) / 240 }},
				&CallbackMatrixObserver{Callback: func(i, j int) float64 { return math.Sin(0.1 * float64(i)) }},
				&CallbackMatrixObserver{Callback: func(i, j int) float64 { return float64(j) / 160 }},
			),
			Layers: []int{0, 1, 3},
		}))
	app.AddSystem(&system.FixedTermination{
		Steps: 100,
	})
	assert.Panics(t, app.Run)
}

// Example observer, reporting a matrix filled with a callback(i, j).
type CallbackMatrixObserver struct {
	Callback func(i, j int) float64
	cols     int
	rows     int
	values   []float64
}

func (o *CallbackMatrixObserver) Initialize(w *ecs.World) {
	o.cols = 240
	o.rows = 160
	o.values = make([]float64, o.cols*o.rows)
}

func (o *CallbackMatrixObserver) Update(w *ecs.World) {}

func (o *CallbackMatrixObserver) Dims() (int, int) {
	return o.cols, o.rows
}

func (o *CallbackMatrixObserver) Values(w *ecs.World) []float64 {
	for idx := 0; idx < len(o.values); idx++ {
		i := idx % o.cols
		j := idx / o.cols
		o.values[idx] = o.Callback(i, j)
	}
	return o.values
}
