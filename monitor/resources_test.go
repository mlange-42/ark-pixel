package monitor_test

import (
	"testing"

	"github.com/mlange-42/ark-pixel/monitor"
	"github.com/mlange-42/ark-pixel/window"
	"github.com/mlange-42/ark-tools/app"
	"github.com/mlange-42/ark-tools/system"
	"github.com/mlange-42/ark/ecs"
)

func ExampleResources() {
	// Create a new model.
	app := app.New()

	// Limit the the simulation speed.
	app.TPS = 30

	// Create a window with a Resources drawer.
	app.AddUISystem((&window.Window{}).
		With(&monitor.Resources{}))

	// Add a termination system that ends the simulation.
	app.AddSystem(&system.FixedTermination{
		Steps: 100,
	})

	app.Run()

	// Run the simulation.
	// Due to the use of the OpenGL UI system, the model must be run via [window.Run].
	// Comment out the code line above, and uncomment the next line to run this example stand-alone.

	// window.Run(app)

	// Output:
}

func TestResources(t *testing.T) {
	app := app.New()

	_ = ecs.AddResource[Position](&app.World, &Position{})
	_ = ecs.ResourceID[Velocity](&app.World)

	app.TPS = 30

	app.AddUISystem((&window.Window{}).
		With(&monitor.Resources{}))

	app.AddSystem(&system.FixedTermination{
		Steps: 10,
	})

	app.Run()
}
