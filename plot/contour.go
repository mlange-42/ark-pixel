package plot

import (
	"fmt"
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

// Contour plot drawer.
//
// Plots a grid as a contours.
// For large grids, this is relatively slow.
// Consider using [Image] instead.
type Contour struct {
	Observer   observer.Grid   // Observers providing a Grid for contours.
	Levels     []float64       // Levels for iso lines. Optional.
	Palette    palette.Palette // Color palette. Optional.
	Labels     Labels          // Labels for plot and axes. Optional.
	HideLegend bool            // Hides the legend.

	data  plotGrid
	scale float64
}

// Initialize the drawer.
func (c *Contour) Initialize(w *ecs.World, win *opengl.Window) {
	c.Observer.Initialize(w)
	c.data = plotGrid{
		Grid: c.Observer,
	}
	c.scale = calcScaleCorrection()
}

// Update the drawer.
func (c *Contour) Update(w *ecs.World) {
	c.Observer.Update(w)
}

// UpdateInputs handles input events of the previous frame update.
func (c *Contour) UpdateInputs(w *ecs.World, win *opengl.Window) {}

// Draw the drawer.
func (c *Contour) Draw(w *ecs.World, win *opengl.Window) {
	width := win.Canvas().Bounds().W()
	height := win.Canvas().Bounds().H()

	c.updateData(w)

	canvas := vgimg.New(vg.Points(width*c.scale)-10, vg.Points(height*c.scale)-10)

	p := plot.New()
	setLabels(p, c.Labels)

	p.X.Tick.Marker = removeLastTicks{}

	cols := c.Palette.Colors()
	min := 0.0
	max := 1.0
	if len(c.Levels) > 0 {
		min = c.Levels[0]
		max = c.Levels[len(c.Levels)-1]
	} else {
		c.Levels = []float64{0.01, 0.05, 0.25, 0.5, 0.75, 0.95, 0.99}
	}

	contours := plotter.Contour{
		GridXYZ:    &c.data,
		Levels:     c.Levels,
		LineStyles: []draw.LineStyle{plotter.DefaultLineStyle},
		Palette:    c.Palette,
		Underflow:  cols[0],
		Overflow:   cols[len(cols)-1],
		Min:        min,
		Max:        max,
	}

	if !c.HideLegend {
		p.Legend = plot.NewLegend()
		p.Legend.TextStyle.Font.Variant = "Mono"
		c.populateLegend(&p.Legend, &contours)
	}

	p.Add(&contours)

	win.Clear(color.White)
	p.Draw(draw.New(canvas))

	img := canvas.Image()
	picture := pixel.PictureDataFromImage(img)

	sprite := pixel.NewSprite(picture, picture.Bounds())
	sprite.Draw(win, pixel.IM.Moved(pixel.V(picture.Rect.W()/2.0+5, picture.Rect.H()/2.0+5)))
}

func (c *Contour) updateData(w *ecs.World) {
	c.data.Values = c.Observer.Values(w)
}

func (c *Contour) populateLegend(legend *plot.Legend, contours *plotter.Contour) {
	var pal []color.Color
	if c.Palette != nil {
		pal = c.Palette.Colors()
	}
	ps := float64(len(pal)-1) / (c.Levels[len(c.Levels)-1] - c.Levels[0])
	if len(c.Levels) == 1 {
		ps = 0
	}
	for i := len(c.Levels) - 1; i >= 0; i-- {
		z := c.Levels[i]
		var col color.Color
		switch {
		case z < contours.Min:
			col = contours.Underflow
		case z > contours.Max:
			col = contours.Overflow
		case len(pal) == 0:
			col = contours.Underflow
		default:
			col = pal[int((z-c.Levels[0])*ps+0.5)] // Apply palette scaling.
		}
		legend.Add(fmt.Sprintf("%f", z), colorThumbnailer{col})
	}
}

// colorThumbnailer implements the Thumbnailer interface.
type colorThumbnailer struct {
	color color.Color
}

// Thumbnail satisfies the plot.Thumbnailer interface.
func (t colorThumbnailer) Thumbnail(c *draw.Canvas) {
	pts := []vg.Point{
		{X: c.Min.X, Y: c.Min.Y},
		{X: c.Min.X, Y: c.Max.Y},
		{X: c.Max.X, Y: c.Max.Y},
		{X: c.Max.X, Y: c.Min.Y},
	}
	poly := c.ClipPolygonY(pts)
	c.FillPolygon(t.color, poly)
}
