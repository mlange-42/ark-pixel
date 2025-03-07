package monitor

import (
	"fmt"
	"image/color"
	"math"

	px "github.com/gopxl/pixel/v2"
	"github.com/gopxl/pixel/v2/backends/opengl"
	"github.com/gopxl/pixel/v2/ext/imdraw"
	"github.com/gopxl/pixel/v2/ext/text"
	"github.com/mlange-42/ark-tools/app"
	"github.com/mlange-42/ark/ecs"
)

// Controls adds UI and keyboard input for controlling the simulation.
// UI controls are displayed in the bottom right corner of the window.
//
// Pause and resume the simulation via a button or by pressing SPACE.
// Manipulate simulation speed (TPS) using buttons or UP/DOWN keys.
//
// Expects a world resource of type Systems ([github.com/mlange-42/ark-tools/model.Systems]).
type Controls struct {
	Scale      float64 // Spatial scaling: cell size in screen pixels. Optional, default 1.
	drawer     imdraw.IMDraw
	systemsRes ecs.Resource[app.Systems]
	text       *text.Text
}

// Initialize the system
func (c *Controls) Initialize(w *ecs.World, win *opengl.Window) {
	c.systemsRes = ecs.NewResource[app.Systems](w)
	if !c.systemsRes.Has() {
		panic("resource of type Systems expected in Controls drawer")
	}

	if c.Scale <= 0 {
		c.Scale = 1
	}

	c.drawer = *imdraw.New(nil)
	c.text = text.New(px.V(0, 0), defaultFont)

}

// Update the drawer.
func (c *Controls) Update(w *ecs.World) {}

// UpdateInputs handles input events of the previous frame update.
func (c *Controls) UpdateInputs(w *ecs.World, win *opengl.Window) {
	sys := c.systemsRes.Get()
	if win.JustPressed(px.KeySpace) {
		sys.Paused = !sys.Paused
		return
	}
	if win.JustPressed(px.KeyUp) {
		sys.TPS = calcTps(sys.TPS, true)
		return
	}
	if win.JustPressed(px.KeyDown) {
		sys.TPS = calcTps(sys.TPS, false)
		return
	}

	if win.JustPressed(px.MouseButton1) {
		width := win.Canvas().Bounds().W()
		height := win.Canvas().Bounds().H()

		mouse := win.MousePosition()
		if c.pauseBounds(width, height).Contains(mouse.X, mouse.Y) {
			sys.Paused = !sys.Paused
		} else if c.upButton(width, height).Contains(mouse.X, mouse.Y) {
			sys.TPS = calcTps(sys.TPS, true)
		} else if c.downButton(width, height).Contains(mouse.X, mouse.Y) {
			sys.TPS = calcTps(sys.TPS, false)
		}
	}
}

// Draw the system
func (c *Controls) Draw(w *ecs.World, win *opengl.Window) {
	width := win.Canvas().Bounds().W()
	height := win.Canvas().Bounds().H()

	sys := c.systemsRes.Get()
	text := "Pause"
	if sys.Paused {
		text = "Resume"
	}
	c.drawButton(c.pauseBounds(width, height), text, win)

	c.drawButton(c.upButton(width, height), "+", win)
	c.drawButton(c.downButton(width, height), "-", win)
	c.drawButton(c.tpsButton(width, height), fmt.Sprintf("%.0f TPS", sys.TPS), win)
}

func (c *Controls) drawButton(b *button, text string, win *opengl.Window) {
	dr := &c.drawer

	dr.Color = color.Black
	dr.Push(px.V(b.X, b.Y), px.V(b.X+b.W, b.Y+b.H))
	dr.Rectangle(0)
	dr.Reset()

	dr.Color = color.White
	dr.Push(px.V(b.X, b.Y), px.V(b.X+b.W, b.Y+b.H))
	dr.Rectangle(1)
	dr.Reset()

	dr.Draw(win)
	dr.Clear()

	c.text.Clear()
	fmt.Fprint(c.text, text)

	wTxt := c.text.Bounds().W()
	hTxt := c.text.Bounds().H()
	cx, cy := b.Center()
	c.text.Draw(win, px.IM.Scaled(px.V(wTxt/2, hTxt/2), c.Scale).Moved(px.V(math.Floor(cx-wTxt/2), math.Floor(cy-hTxt/2))))
}

func (c *Controls) pauseBounds(w, h float64) *button {
	return &button{
		w - 85*c.Scale,
		5 + 15*c.Scale,
		60 * c.Scale,
		30 * c.Scale,
	}
}

func (c *Controls) upButton(w, h float64) *button {
	return &button{
		w - 25*c.Scale,
		5 + 30*c.Scale,
		20 * c.Scale,
		15 * c.Scale,
	}
}

func (c *Controls) downButton(w, h float64) *button {
	return &button{
		w - 25*c.Scale,
		5 + 15*c.Scale,
		20 * c.Scale,
		15 * c.Scale,
	}
}

func (c *Controls) tpsButton(w, h float64) *button {
	return &button{
		w - 85*c.Scale,
		5,
		80 * c.Scale,
		15 * c.Scale,
	}
}

type button struct {
	X float64
	Y float64
	W float64
	H float64
}

func (b *button) Center() (float64, float64) {
	return b.X + b.W/2, b.Y + b.H/2
}

func (b *button) Contains(x, y float64) bool {
	return x >= b.X && y >= b.Y && x <= b.X+b.W && y <= b.Y+b.H
}
