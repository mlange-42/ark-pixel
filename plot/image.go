package plot

import (
	"image/color"

	pixel "github.com/gopxl/pixel/v2"
	"github.com/gopxl/pixel/v2/backends/opengl"
	"github.com/mazznoer/colorgrad"
	"github.com/mlange-42/ark-pixel/window"
	"github.com/mlange-42/ark-tools/observer"
	"github.com/mlange-42/ark/ecs"
)

// Image drawer.
//
// Draws an image from a Matrix observer.
// The image is scaled to the canvas extent, with preserved aspect ratio.
// Does not add plot axes etc.
type Image struct {
	Scale    float64            // Spatial scaling: cell size in screen pixels. Optional, default auto.
	Observer observer.Matrix    // Observer providing 2D matrix or grid data.
	Colors   colorgrad.Gradient // Colors for mapping values.
	Min      float64            // Minimum value for color mapping. Optional.
	Max      float64            // Maximum value for color mapping. Optional. Is set to 1.0 if both Min and Max are zero.
	slope    float64
	picture  *pixel.PictureData
}

// Initialize the system
func (i *Image) Initialize(w *ecs.World, win *opengl.Window) {
	i.Observer.Initialize(w)

	if i.Min == 0 && i.Max == 0 {
		i.Max = 1
	}

	i.slope = 1.0 / (i.Max - i.Min)

	width, height := i.Observer.Dims()
	i.picture = pixel.MakePictureData(pixel.R(0, 0, float64(width), float64(height)))
}

// Update the drawer.
func (i *Image) Update(w *ecs.World) {
	i.Observer.Update(w)
}

// UpdateInputs handles input events of the previous frame update.
func (i *Image) UpdateInputs(w *ecs.World, win *opengl.Window) {}

// Draw the system
func (i *Image) Draw(w *ecs.World, win *opengl.Window) {
	values := i.Observer.Values(w)

	length := len(values)
	for j := 0; j < length; j++ {
		i.picture.Pix[j] = i.valueToColor(values[j])
	}

	scale := i.Scale
	if i.Scale <= 0 {
		scale = window.Scale(win, i.picture.Rect.W(), i.picture.Rect.H())
	}

	sprite := pixel.NewSprite(i.picture, i.picture.Bounds())
	sprite.Draw(win,
		pixel.IM.Moved(pixel.V(i.picture.Rect.W()/2.0, i.picture.Rect.H()/2.0)).
			Scaled(pixel.Vec{}, scale),
	)
}

func (i *Image) valueToColor(v float64) color.RGBA {
	c := i.Colors.At((v - i.Min) * i.slope)
	return color.RGBA{
		R: uint8(c.R * 255),
		G: uint8(c.G * 255),
		B: uint8(c.B * 255),
		A: 0xff,
	}
}
