package monitor_test

import (
	"github.com/mlange-42/ark-pixel/monitor"
	"github.com/mlange-42/ark-pixel/window"
	"github.com/mlange-42/ark-tools/app"
	"github.com/mlange-42/ark-tools/system"
)

func ExampleControls() {
	// Create a new model.
	app := app.New()

	// Limit the the simulation speed.
	app.TPS = 30

	// Create a window with a Controls drawer.
	app.AddUISystem((&window.Window{}).
		With(&monitor.Controls{Scale: 2}))

	// Controls is intended as an overlay, so more drawers can be added before it.
	app.AddUISystem((&window.Window{}).
		With(
			&monitor.Monitor{},
			&monitor.Controls{},
		))

	// Add a termination system that ends the simulation.
	app.AddSystem(&system.FixedTermination{
		Steps: 100,
	})

	app.Run()

	// Run the simulation.
	// Due to the use of the OpenGL UI system, the model must be run via [window.Run].
	// Comment out the code line above, and uncomment the next line to run this example stand-alone.

	//window.Run(app)

	// Output:
}
