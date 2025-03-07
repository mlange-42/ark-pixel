package plot

import (
	"fmt"
	"image/color"

	pixel "github.com/gopxl/pixel/v2"
	"github.com/gopxl/pixel/v2/backends/opengl"
	"github.com/mlange-42/ark-pixel/window"
	"github.com/mlange-42/ark-tools/observer"
	"github.com/mlange-42/ark/ecs"
)

// ImageRGB drawer.
//
// Draws an image from a Matrix observer per RGB color channel.
// The image is scaled to the canvas extent, with preserved aspect ratio.
// Does not add plot axes etc.
type ImageRGB struct {
	Scale    float64               // Spatial scaling: cell size in screen pixels. Optional, default auto.
	Observer observer.MatrixLayers // Observer providing data for color channels.
	Layers   []int                 // Layer indices. Optional, defaults to [0, 1, 2]. Use -1 to ignore a channel.
	Min      []float64             // Minimum value for channel color mapping. Optional, default [0, 0, 0].
	Max      []float64             // Maximum value for channel color mapping. Optional, default [1, 1, 1].
	slope    []float64
	dataLen  int
	picture  *pixel.PictureData
}

// Initialize the drawer.
func (i *ImageRGB) Initialize(w *ecs.World, win *opengl.Window) {
	i.Observer.Initialize(w)

	if i.Layers == nil {
		i.Layers = []int{0, 1, 2}
	} else if len(i.Layers) != 3 {
		panic("rgb image plot Layers must be of length 3")
	}

	layers := i.Observer.Layers()
	for _, l := range i.Layers {
		if l >= 0 && layers <= l {
			panic(fmt.Sprintf("layer index %d out of range", l))
		}
	}

	if i.Min == nil {
		i.Min = []float64{0, 0, 0}
	}
	if i.Max == nil {
		i.Max = []float64{1, 1, 1}
	}
	if len(i.Min) != 3 {
		panic("RgbImage plot needs exactly 3 Min values")
	}
	if len(i.Max) != 3 {
		panic("RgbImage plot needs exactly 3 Max values")
	}

	i.slope = []float64{
		1.0 / (i.Max[0] - i.Min[0]),
		1.0 / (i.Max[1] - i.Min[1]),
		1.0 / (i.Max[2] - i.Min[2]),
	}

	width, height := i.Observer.Dims()
	i.dataLen = width * height
	i.picture = pixel.MakePictureData(pixel.R(0, 0, float64(width), float64(height)))
}

// Update the drawer.
func (i *ImageRGB) Update(w *ecs.World) {
	i.Observer.Update(w)
}

// UpdateInputs handles input events of the previous frame update.
func (i *ImageRGB) UpdateInputs(w *ecs.World, win *opengl.Window) {}

// Draw the drawer.
func (i *ImageRGB) Draw(w *ecs.World, win *opengl.Window) {
	cannels := i.Observer.Values(w)

	values := append([]float64{}, i.Min...)
	for j := 0; j < i.dataLen; j++ {
		for i, k := range i.Layers {
			if k >= 0 {
				values[i] = cannels[k][j]
			}
		}
		i.picture.Pix[j] = i.valuesToColor(values[0], values[1], values[2])
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

func (i *ImageRGB) valuesToColor(r, g, b float64) color.RGBA {
	return color.RGBA{
		R: norm(r, i.Min[0], i.slope[0]),
		G: norm(g, i.Min[1], i.slope[1]),
		B: norm(b, i.Min[2], i.slope[2]),
		A: 0xff,
	}
}

func norm(v, off, slope float64) uint8 {
	vv := (v - off) * slope
	if vv <= 0 {
		return 0
	}
	if vv >= 1 {
		return 255
	}
	return uint8(vv * 255)
}
