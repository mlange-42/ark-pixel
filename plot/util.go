package plot

import (
	"fmt"
	"image/color"

	"github.com/mlange-42/ark-tools/observer"
	"golang.org/x/image/colornames"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/vg/vgimg"
)

var preferredTps = []float64{0, 1, 2, 3, 4, 5, 7, 10, 15, 20, 30, 40, 50, 60, 80, 100, 120, 150, 200, 250, 500, 750, 1000, 2000, 5000, 10000}

var defaultColors = []color.Color{
	colornames.Blue,
	colornames.Orange,
	colornames.Green,
	colornames.Purple,
	colornames.Red,
	colornames.Turquoise,
}

// Labels for plots.
type Labels struct {
	Title string // Plot title
	X     string // X axis label
	Y     string // Y axis label
}

// Get the index of an element in a slice.
func find[T comparable](sl []T, value T) (int, bool) {
	for i, v := range sl {
		if v == value {
			return i, true
		}
	}
	return -1, false
}

// Calculate scale correction for scaled monitors.
func calcScaleCorrection() float64 {
	return 72.0 / vgimg.DefaultDPI
}

func setLabels(p *plot.Plot, l Labels) {
	p.Title.Text = l.Title
	p.Title.TextStyle.Font.Size = 16
	p.Title.TextStyle.Font.Variant = "Mono"

	p.X.Label.Text = l.X
	p.X.Label.TextStyle.Font.Size = 14
	p.X.Label.TextStyle.Font.Variant = "Mono"

	p.X.Tick.Label.Font.Size = 12
	p.X.Tick.Label.Font.Variant = "Mono"

	p.Y.Label.Text = l.Y
	p.Y.Label.TextStyle.Font.Size = 14
	p.Y.Label.TextStyle.Font.Variant = "Mono"

	p.Y.Tick.Label.Font.Size = 12
	p.Y.Tick.Label.Font.Variant = "Mono"

	p.Y.Tick.Marker = paddedTicks{}
}

// Left-pads tick labels to avoid jumping Y axis.
type paddedTicks struct {
	plot.DefaultTicks
}

func (t paddedTicks) Ticks(min, max float64) []plot.Tick {
	ticks := t.DefaultTicks.Ticks(min, max)
	for i := 0; i < len(ticks); i++ {
		ticks[i].Label = fmt.Sprintf("%*s", 10, ticks[i].Label)
	}
	return ticks
}

// Removes the last tick label to avoid jumping X axis.
type removeLastTicks struct {
	plot.DefaultTicks
}

func (t removeLastTicks) Ticks(min, max float64) []plot.Tick {
	ticks := t.DefaultTicks.Ticks(min, max)
	for i := 0; i < len(ticks); i++ {
		tick := &ticks[i]
		if tick.IsMinor() {
			continue
		}
		if tick.Value > max-(0.05*(max-min)) {
			tick.Label = ""
		}
	}
	return ticks
}

type plotGrid struct {
	observer.Grid
	Values []float64
}

func (g *plotGrid) Z(c, r int) float64 {
	w, _ := g.Dims()
	return g.Values[r*w+c]
}

type ringBuffer[T any] struct {
	data  []T
	start int
}

func newRingBuffer[T any](cap int) ringBuffer[T] {
	return ringBuffer[T]{
		data:  make([]T, 0, cap),
		start: 0,
	}
}

func (r *ringBuffer[T]) Len() int {
	return len(r.data)
}

func (r *ringBuffer[T]) Cap() int {
	return cap(r.data)
}

func (r *ringBuffer[T]) Get(idx int) T {
	return r.data[(r.start+idx)%r.Cap()]
}

func (r *ringBuffer[T]) Add(elem T) {
	if cap(r.data) > len(r.data) {
		r.data = append(r.data, elem)
		return
	}
	r.data[r.start] = elem
	r.start = (r.start + 1) % r.Cap()
}
