package plot

import (
	"image/color"

	pixel "github.com/gopxl/pixel/v2"
	"github.com/gopxl/pixel/v2/backends/opengl"
	"github.com/mlange-42/ark-tools/observer"
	"github.com/mlange-42/ark/ecs"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/palette"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/vg/vgimg"
)

// HeatMap plot drawer.
//
// Plots a grid as a heatmap image.
// For large grids, this is relatively slow.
// Consider using [Image] instead.
type HeatMap struct {
	Observer observer.Grid   // Observers providing a Grid for contours.
	Palette  palette.Palette // Color palette. Optional.
	Min      float64         // Minimum value for color mapping. Optional.
	Max      float64         // Maximum value for color mapping. Optional. Is set to 1.0 if both Min and Max are zero.
	Labels   Labels          // Labels for plot and axes. Optional.

	data  plotGrid
	scale float64
}

// Initialize the drawer.
func (h *HeatMap) Initialize(w *ecs.World, win *opengl.Window) {
	h.Observer.Initialize(w)
	h.data = plotGrid{
		Grid: h.Observer,
	}

	h.scale = calcScaleCorrection()

	if h.Min == 0 && h.Max == 0 {
		h.Max = 1
	}
}

// Update the drawer.
func (h *HeatMap) Update(w *ecs.World) {
	h.Observer.Update(w)
}

// UpdateInputs handles input events of the previous frame update.
func (h *HeatMap) UpdateInputs(w *ecs.World, win *opengl.Window) {}

// Draw the drawer.
func (h *HeatMap) Draw(w *ecs.World, win *opengl.Window) {
	width := win.Canvas().Bounds().W()
	height := win.Canvas().Bounds().H()

	h.updateData(w)

	c := vgimg.New(vg.Points(width*h.scale)-10, vg.Points(height*h.scale)-10)

	p := plot.New()
	setLabels(p, h.Labels)

	p.X.Tick.Marker = removeLastTicks{}

	cols := h.Palette.Colors()
	heat := plotter.HeatMap{
		GridXYZ:    &h.data,
		Palette:    h.Palette,
		Rasterized: false,
		Underflow:  cols[0],
		Overflow:   cols[len(cols)-1],
		Min:        h.Min,
		Max:        h.Max,
	}

	p.Add(&heat)

	win.Clear(color.White)
	p.Draw(draw.New(c))

	img := c.Image()
	picture := pixel.PictureDataFromImage(img)

	sprite := pixel.NewSprite(picture, picture.Bounds())
	sprite.Draw(win, pixel.IM.Moved(pixel.V(picture.Rect.W()/2.0+5, picture.Rect.H()/2.0+5)))
}

func (h *HeatMap) updateData(w *ecs.World) {
	h.data.Values = h.Observer.Values(w)
}
