package monitor

import (
	"math"

	"github.com/gopxl/pixel/v2/ext/text"
	"golang.org/x/image/font/basicfont"
	"gonum.org/v1/plot/vg/vgimg"
)

var defaultFont = text.NewAtlas(basicfont.Face7x13, text.ASCII)

var preferredTicks = []float64{1, 2, 5, 10}
var preferredTps = []float64{0, 1, 2, 3, 4, 5, 7, 10, 15, 20, 30, 40, 50, 60, 80, 100, 120, 150, 200, 250, 500, 750, 1000, 2000, 5000, 10000}

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

// Calculate the optimal step size for axis ticks.
func calcTicksStep(max float64, desired int) float64 {
	steps := float64(desired)
	approxStep := float64(max) / (steps - 1)
	stepPower := math.Pow(10, -math.Floor(math.Log10(approxStep)))
	normalizedStep := approxStep * stepPower
	for _, s := range preferredTicks {
		if s >= normalizedStep {
			normalizedStep = s
			break
		}
	}
	return normalizedStep / stepPower
}

// Calculate TPS when increasing/decreasing it.
func calcTps(curr float64, increase bool) float64 {
	ln := len(preferredTps)
	if increase {
		for i := 0; i < ln; i++ {
			if preferredTps[i] > curr {
				return preferredTps[i]
			}
		}
		return curr
	}
	for i := 1; i < ln; i++ {
		if preferredTps[i] >= curr {
			return preferredTps[i-1]
		}
	}
	return 0
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
