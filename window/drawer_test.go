package window_test

import (
	"image/color"

	pixel "github.com/gopxl/pixel/v2"
	"github.com/gopxl/pixel/v2/backends/opengl"
	"github.com/gopxl/pixel/v2/ext/imdraw"
	"github.com/mlange-42/ark-pixel/window"
	"github.com/mlange-42/ark-tools/resource"
	"github.com/mlange-42/ark/ecs"
)

func ExampleDrawer() {
	var dr window.Drawer = &RectDrawer{}
	_ = dr
	// Output:
}

// RectDrawer is an example drawer.
type RectDrawer struct {
	dr imdraw.IMDraw
}

// Initialize the RectDrawer.
func (d *RectDrawer) Initialize(w *ecs.World, win *opengl.Window) {
	// Create a drawer from the Pixel engine.
	d.dr = *imdraw.New(nil)
}

// Update the RectDrawer (does nothing).
func (d *RectDrawer) Update(w *ecs.World) {}

// UpdateInputs handles input events of the previous frame update.
func (d *RectDrawer) UpdateInputs(w *ecs.World, win *opengl.Window) {}

// Draw the RectDrawer's stuff.
func (d *RectDrawer) Draw(w *ecs.World, win *opengl.Window) {
	// Get a resource from the world.
	tick := ecs.GetResource[resource.Tick](w)
	offset := float64(tick.Tick)

	// Create a white rectangle that moves with progressing model tick.
	d.dr.Color = color.White
	d.dr.Push(pixel.V(50+offset, 50+offset), pixel.V(250+offset, 200+offset))
	d.dr.Rectangle(0)

	// Draw everything on the window.
	d.dr.Draw(win)

	// Reset the drawer
	d.dr.Reset()
	d.dr.Clear()
}
