package plot

import (
	"fmt"
	"image/color"

	pixel "github.com/gopxl/pixel/v2"
	"github.com/gopxl/pixel/v2/backends/opengl"
	"github.com/mlange-42/ark-tools/observer"
	"github.com/mlange-42/ark/ecs"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/vg/vgimg"
)

// TimeSeries plot drawer.
//
// Creates a line series per column of the observer.
// Adds one row to the data per update.
type TimeSeries struct {
	Observer       observer.Row // Observer providing a data row per update.
	Columns        []string     // Columns to show, by name. Optional, default all.
	UpdateInterval int          // Interval for getting data from the the observer, in model ticks. Optional.
	Labels         Labels       // Labels for plot and axes. Optional.
	MaxRows        int          // Maximum number of rows to keep. Zero means unlimited. Optional.

	indices []int
	headers []string
	series  []plotter.XYs
	scale   float64
	step    int64
}

// append a y value to each series, with a common x value.
func (t *TimeSeries) append(x float64, values []float64) {
	for i := 0; i < len(t.series); i++ {
		t.series[i] = append(t.series[i], plotter.XY{X: x, Y: values[i]})
		if t.MaxRows > 0 && len(t.series[i]) > t.MaxRows {
			t.series[i] = t.series[i][len(t.series[i])-t.MaxRows:]
		}
	}
}

// Initialize the drawer.
func (t *TimeSeries) Initialize(w *ecs.World, win *opengl.Window) {
	t.Observer.Initialize(w)

	t.headers = t.Observer.Header()

	if len(t.Columns) == 0 {
		t.indices = make([]int, len(t.headers))
		for i := 0; i < len(t.indices); i++ {
			t.indices[i] = i
		}
	} else {
		t.indices = make([]int, len(t.Columns))
		var ok bool
		for i := 0; i < len(t.indices); i++ {
			t.indices[i], ok = find(t.headers, t.Columns[i])
			if !ok {
				panic(fmt.Sprintf("column '%s' not found", t.Columns[i]))
			}
		}
	}

	t.series = make([]plotter.XYs, len(t.headers))

	t.scale = calcScaleCorrection()
	t.step = 0
}

// Update the drawer.
func (t *TimeSeries) Update(w *ecs.World) {
	t.Observer.Update(w)
	if t.UpdateInterval <= 1 || t.step%int64(t.UpdateInterval) == 0 {
		t.append(float64(t.step), t.Observer.Values(w))
	}
	t.step++
}

// UpdateInputs handles input events of the previous frame update.
func (t *TimeSeries) UpdateInputs(w *ecs.World, win *opengl.Window) {}

// Draw the drawer.
func (t *TimeSeries) Draw(w *ecs.World, win *opengl.Window) {
	width := win.Canvas().Bounds().W()
	height := win.Canvas().Bounds().H()

	c := vgimg.New(vg.Points(width*t.scale)-10, vg.Points(height*t.scale)-10)

	p := plot.New()
	setLabels(p, t.Labels)

	p.X.Tick.Marker = removeLastTicks{}

	p.Legend = plot.NewLegend()
	p.Legend.TextStyle.Font.Variant = "Mono"

	for i, idx := range t.indices {
		lines, err := plotter.NewLine(t.series[idx])
		if err != nil {
			panic(err)
		}
		lines.Color = defaultColors[i%len(defaultColors)]
		p.Add(lines)
		p.Legend.Add(t.headers[idx], lines)
	}

	win.Clear(color.White)
	p.Draw(draw.New(c))

	img := c.Image()
	picture := pixel.PictureDataFromImage(img)

	sprite := pixel.NewSprite(picture, picture.Bounds())
	sprite.Draw(win, pixel.IM.Moved(pixel.V(picture.Rect.W()/2.0+5, picture.Rect.H()/2.0+5)))
}
