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

// Field plot drawer.
//
// Plots a vector field from a GridLayers observer.
// For large grids, this is relatively slow.
// Consider using [ImageRGB] instead.
type Field struct {
	Observer observer.GridLayers // Observers providing field component grids.
	Labels   Labels              // Labels for plot and axes. Optional.
	Layers   []int               // Layer indices. Optional, defaults to (0, 1).

	data  plotField
	scale float64
}

// Initialize the drawer.
func (f *Field) Initialize(w *ecs.World, win *opengl.Window) {
	f.Observer.Initialize(w)

	f.data = plotField{
		GridLayers: f.Observer,
	}

	if f.Layers == nil {
		f.Layers = []int{0, 1}
	} else if len(f.Layers) != 2 {
		panic("field plot Layers must be of length 2")
	}
	layers := f.Observer.Layers()
	for _, l := range f.Layers {
		if layers <= l {
			panic(fmt.Sprintf("layer index %d out of range", l))
		}
	}

	f.scale = calcScaleCorrection()
}

// Update the drawer.
func (f *Field) Update(w *ecs.World) {
	f.Observer.Update(w)
}

// UpdateInputs handles input events of the previous frame update.
func (f *Field) UpdateInputs(w *ecs.World, win *opengl.Window) {}

// Draw the drawer.
func (f *Field) Draw(w *ecs.World, win *opengl.Window) {
	width := win.Canvas().Bounds().W()
	height := win.Canvas().Bounds().H()

	f.updateData(w)

	canvas := vgimg.New(vg.Points(width*f.scale)-10, vg.Points(height*f.scale)-10)

	p := plot.New()
	setLabels(p, f.Labels)

	p.X.Tick.Marker = removeLastTicks{}

	field := plotter.NewField(&f.data)

	p.Add(field)

	win.Clear(color.White)
	p.Draw(draw.New(canvas))

	img := canvas.Image()
	picture := pixel.PictureDataFromImage(img)

	sprite := pixel.NewSprite(picture, picture.Bounds())
	sprite.Draw(win, pixel.IM.Moved(pixel.V(picture.Rect.W()/2.0+5, picture.Rect.H()/2.0+5)))
}

func (f *Field) updateData(w *ecs.World) {
	values := f.Observer.Values(w)
	f.data.XValues = values[f.Layers[0]]
	f.data.YValues = values[f.Layers[1]]
}

type plotField struct {
	observer.GridLayers
	XValues []float64
	YValues []float64
}

func (f *plotField) Vector(c, r int) plotter.XY {
	w, _ := f.GridLayers.Dims()
	return plotter.XY{
		X: f.XValues[r*w+c],
		Y: f.YValues[r*w+c],
	}
}
