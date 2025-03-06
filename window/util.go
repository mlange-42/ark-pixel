package window

import (
	"math"

	"github.com/gopxl/pixel/v2/backends/opengl"
)

// Scale calculates the drawing scale for fitting a source region into a window's canvas.
func Scale(win *opengl.Window, srcWidth, srcHeight float64) float64 {
	winWidth, winHeight := win.Canvas().Bounds().W(), win.Canvas().Bounds().H()
	scX, scY := winWidth/float64(srcWidth), winHeight/float64(srcHeight)
	return math.Min(scX, scY)
}
