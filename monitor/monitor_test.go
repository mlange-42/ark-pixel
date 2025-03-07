package monitor_test

import (
	"github.com/mlange-42/ark-pixel/monitor"
	"github.com/mlange-42/ark-pixel/window"
	"github.com/mlange-42/ark-tools/app"
	"github.com/mlange-42/ark-tools/system"
)

func ExampleMonitor() {
	// Create a new model.
	app := app.New()

	// Limit the the simulation speed.
	app.TPS = 30

	// Create a window with a Monitor drawer.
	app.AddUISystem((&window.Window{}).
		With(&monitor.Monitor{}))

	// Add a termination system that ends the simulation.
	app.AddSystem(&system.FixedTermination{
		Steps: 100,
	})

	// Run the simulation.
	// Due to the use of the OpenGL UI system, the model must be run via [window.Run].
	// Uncomment the next line. It is commented out as the CI has no display device to test the model run.

	//window.Run(app)

	// Output:
}

func ExampleNewMonitorWindow() {
	// Create a new model.
	app := app.New()

	// Limit the the simulation speed.
	app.TPS = 30

	// Create a window with a Monitor drawer, using the shorthand constructor.
	app.AddUISystem(monitor.NewMonitorWindow(10))

	// Add a termination system that ends the simulation.
	app.AddSystem(&system.FixedTermination{
		Steps: 100,
	})

	// Run the simulation.
	// Due to the use of the OpenGL UI system, the model must be run via [window.Run].
	// Uncomment the next line. It is commented out as the CI has no display device to test the model run.

	//window.Run(app)

	// Output:
}
