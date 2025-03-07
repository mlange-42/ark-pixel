package window

import (
	"fmt"
	"log"

	pixel "github.com/gopxl/pixel/v2"
	"github.com/gopxl/pixel/v2/backends/opengl"
	"github.com/mlange-42/ark-tools/resource"
	"github.com/mlange-42/ark/ecs"
	"golang.org/x/image/colornames"
)

// Drawer interface.
// Drawers are used by the [Window] to render information from an Ark application.
type Drawer interface {
	// Initialize is called before any other method.
	// Use it to initialize the Drawer.
	Initialize(w *ecs.World, win *opengl.Window)

	// Update is called with normal system updates.
	// Can be used to update observers.
	Update(w *ecs.World)

	// UpdateInputs is called on every UI update, i.e. with the frequency of FPS.
	// Can be used to handle user input of the previous frame update.
	UpdateInputs(w *ecs.World, win *opengl.Window)

	// Draw is called on UI updates, every [Model.DrawInterval] steps.
	// Draw is not called when the host window is minimized.
	// Do all OpenGL drawing on the window here.
	Draw(w *ecs.World, win *opengl.Window)
}

// Bounds define a bounding box for a window.
type Bounds struct {
	X int // X position
	Y int // Y position
	W int // Width
	H int // Height
}

// B created a new Bounds object.
func B(x, y, w, h int) Bounds {
	return Bounds{x, y, w, h}
}

// Window provides an OpenGL window for drawing.
// Drawing is done by one or more [Drawer] instances.
// Further, window bounds and update and draw intervals can be configured.
//
// If the world contains a resource of type [github.com/mlange-42/ark-tools/resource/Termination],
// the model is terminated when the window is closed.
type Window struct {
	Title        string   // Window title. Optional.
	Bounds       Bounds   // Window bounds (position and size). Optional.
	Drawers      []Drawer // Drawers in increasing z order.
	DrawInterval int      // Interval for re-drawing, in UI frames. Optional.
	window       *opengl.Window
	drawStep     int64
	isClosed     bool
	termRes      ecs.Resource[resource.Termination]
}

// With adds one or more [Drawer] instances to the window.
func (w *Window) With(drawers ...Drawer) *Window {
	w.Drawers = append(w.Drawers, drawers...)
	return w
}

// Initialize the window system.
func (w *Window) Initialize(world *ecs.World) {}

// InitializeUI the window system.
func (w *Window) InitializeUI(world *ecs.World) {
	if w.Bounds.W <= 0 {
		w.Bounds.W = 1024
	}
	if w.Bounds.H <= 0 {
		w.Bounds.H = 768
	}
	if w.Title == "" {
		w.Title = "Ark"
	}
	cfg := opengl.WindowConfig{
		Title:     w.Title,
		Bounds:    pixel.R(0, 0, float64(w.Bounds.W), float64(w.Bounds.H)),
		Position:  pixel.V(float64(w.Bounds.X), float64(w.Bounds.Y)),
		Resizable: true,
	}

	defer func() {
		if err := recover(); err != nil {
			txt := fmt.Sprint(err)
			if txt == "mainthread: did not call Run" {
				log.Fatal("ERROR: when using graphics via the pixel engine, run the app like this:\n    window.Run(app)")
			}
			panic(err)
		}
	}()

	var err error
	w.window, err = opengl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	for _, d := range w.Drawers {
		d.Initialize(world, w.window)
	}

	w.termRes = ecs.NewResource[resource.Termination](world)
	w.drawStep = 0
	w.isClosed = false
}

// Update the window system.
func (w *Window) Update(world *ecs.World) {
	if w.isClosed {
		return
	}
	for _, d := range w.Drawers {
		d.Update(world)
	}
}

// UpdateUI the window system.
func (w *Window) UpdateUI(world *ecs.World) {
	if w.window.Closed() {
		if !w.isClosed {
			term := w.termRes.Get()
			if term != nil {
				term.Terminate = true
			}
			w.isClosed = true
		}
		return
	}
	if !w.isMinimized() && (w.DrawInterval <= 1 || w.drawStep%int64(w.DrawInterval) == 0) {
		w.window.Clear(colornames.Black)

		for _, d := range w.Drawers {
			d.Draw(world, w.window)
		}
	}
	w.drawStep++
}

func (w *Window) isMinimized() bool {
	b := w.window.Bounds()
	return b.W() <= 0 || b.H() <= 0
}

// PostUpdateUI updates the underlying GL window and input events.
func (w *Window) PostUpdateUI(world *ecs.World) {
	w.window.Update()
	for _, d := range w.Drawers {
		d.UpdateInputs(world, w.window)
	}
}

// Finalize the window system.
func (w *Window) Finalize(world *ecs.World) {}

// FinalizeUI the window system.
func (w *Window) FinalizeUI(world *ecs.World) {
	w.window.Destroy()
}
