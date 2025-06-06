package monitor_test

import (
	"testing"

	"github.com/mlange-42/ark-pixel/monitor"
	"github.com/mlange-42/ark-pixel/window"
	"github.com/mlange-42/ark-tools/app"
	"github.com/mlange-42/ark-tools/system"
)

func ExampleSystems() {
	// Create a new model.
	app := app.New()

	// Limit the the simulation speed.
	app.TPS = 30

	// Create a window with a Systems drawer.
	app.AddUISystem((&window.Window{}).
		With(&monitor.Systems{}))

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

func TestSystems(t *testing.T) {
	app := app.New()
	app.TPS = 300

	app.AddUISystem((&window.Window{}).
		With(&monitor.Systems{
			HideUISystems: true,
			HideNames:     true,
		}))

	app.AddSystem(&system.FixedTermination{
		Steps: 100,
	})

	app.Run()
}
