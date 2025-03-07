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

// Bars plot drawer.
//
// Creates a bar per column of the observer.
type Bars struct {
	Observer observer.Row // Observer providing a data series for bars.
	Columns  []string     // Columns to show, by name. Optional, default all.
	YLim     [2]float64   // Y axis limits. Optional, default auto.
	Labels   Labels       // Labels for plot and axes. Optional.

	indices []int
	headers []string
	series  plotter.Values
	scale   float64
}

// Initialize the drawer.
func (b *Bars) Initialize(w *ecs.World, win *opengl.Window) {
	b.Observer.Initialize(w)

	headers := b.Observer.Header()

	if len(b.Columns) == 0 {
		b.indices = make([]int, len(headers))
		for i := 0; i < len(b.indices); i++ {
			b.indices[i] = i
		}
	} else {
		b.indices = make([]int, len(b.Columns))
		var ok bool
		for i := 0; i < len(b.indices); i++ {
			b.indices[i], ok = find(headers, b.Columns[i])
			if !ok {
				panic(fmt.Sprintf("column '%s' not found", b.Columns[i]))
			}
		}
	}

	b.series = make([]float64, len(b.indices))
	b.headers = make([]string, len(b.indices))
	for i, idx := range b.indices {
		b.headers[i] = headers[idx]
	}

	b.scale = calcScaleCorrection()

}

// Update the drawer.
func (b *Bars) Update(w *ecs.World) {
	b.Observer.Update(w)
}

// UpdateInputs handles input events of the previous frame update.
func (b *Bars) UpdateInputs(w *ecs.World, win *opengl.Window) {}

// Draw the drawer.
func (b *Bars) Draw(w *ecs.World, win *opengl.Window) {
	width := win.Canvas().Bounds().W()
	height := win.Canvas().Bounds().H()

	b.updateData(w)

	c := vgimg.New(vg.Points(width*b.scale)-10, vg.Points(height*b.scale)-10)

	p := plot.New()
	setLabels(p, b.Labels)

	if b.YLim[0] != 0 || b.YLim[1] != 0 {
		p.Y.Min = b.YLim[0]
		p.Y.Max = b.YLim[1]
	}

	bw := 0.5 * (width - 50) / float64(len(b.series))
	bars, err := plotter.NewBarChart(b.series, vg.Points(bw))
	if err != nil {
		panic(err)
	}
	bars.Color = defaultColors[0]
	p.Add(bars)
	p.NominalX(b.headers...)

	win.Clear(color.White)
	p.Draw(draw.New(c))

	img := c.Image()
	picture := pixel.PictureDataFromImage(img)

	sprite := pixel.NewSprite(picture, picture.Bounds())
	sprite.Draw(win, pixel.IM.Moved(pixel.V(picture.Rect.W()/2.0+5, picture.Rect.H()/2.0+5)))
}

func (b *Bars) updateData(w *ecs.World) {
	values := b.Observer.Values(w)

	for i, idx := range b.indices {
		b.series[i] = values[idx]
	}
}
