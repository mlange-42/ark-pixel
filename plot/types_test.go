package plot_test

import (
	"math"
	"math/rand"

	"github.com/mlange-42/ark/ecs"
)

// TableObserver to generate random time series.
type TableObserver struct {
	data [][]float64
}

func (o *TableObserver) Initialize(w *ecs.World) {
	rows := 25
	o.data = make([][]float64, rows)

	for i := 0; i < rows; i++ {
		o.data[i] = []float64{float64(i), float64(i) / float64(rows), float64(rows-i) / float64(rows), 0}
	}
}
func (o *TableObserver) Update(w *ecs.World) {}
func (o *TableObserver) Header() []string {
	return []string{"X", "A", "B", "C"}
}
func (o *TableObserver) Values(w *ecs.World) [][]float64 {
	for i := 0; i < len(o.data); i++ {
		o.data[i][3] = rand.Float64()
	}
	return o.data
}

// TableObserver to generate test time series containing NaN.
type TableObserverNaN struct {
	data [][]float64
}

func (o *TableObserverNaN) Initialize(w *ecs.World) {
	rows := 25
	o.data = make([][]float64, rows)

	for i := 0; i < rows; i++ {
		v := 1.0
		if i < 5 || i > rows-5 {
			v = math.NaN()
		}
		o.data[i] = []float64{float64(i), v}
	}
}
func (o *TableObserverNaN) Update(w *ecs.World) {}
func (o *TableObserverNaN) Header() []string {
	return []string{"X", "A"}
}
func (o *TableObserverNaN) Values(w *ecs.World) [][]float64 {
	return o.data
}
