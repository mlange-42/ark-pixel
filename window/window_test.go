package window_test

import (
	"github.com/mlange-42/ark-pixel/window"
	"github.com/mlange-42/ark-tools/app"
)

func ExampleWindow() {
	app := app.New()

	// Create a Window with at least one Drawer.
	window := (&window.Window{Bounds: window.B(100, 100, 800, 600)}).
		With(&RectDrawer{})

	// Add is to the model as UI system.
	app.AddUISystem(window)
	// Output:
}
