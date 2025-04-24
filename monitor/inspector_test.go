package monitor_test

import (
	"testing"

	"github.com/mlange-42/ark-pixel/monitor"
	"github.com/mlange-42/ark-pixel/window"
	"github.com/mlange-42/ark-tools/app"
	"github.com/mlange-42/ark-tools/resource"
	"github.com/mlange-42/ark-tools/system"
	"github.com/mlange-42/ark/ecs"
)

func ExampleInspector() {
	// Create a new model.
	app := app.New()

	// Limit the the simulation speed.
	app.TPS = 30

	// Create an entity to inspect it.
	posID := ecs.ComponentID[Position](&app.World)
	velID := ecs.ComponentID[Velocity](&app.World)
	entity := app.World.Unsafe().NewEntity(posID, velID)

	// Set it as selected entity.
	ecs.AddResource(&app.World, &resource.SelectedEntity{Selected: entity})

	// Create a window with an Inspector drawer.
	app.AddUISystem((&window.Window{}).
		With(&monitor.Inspector{}))

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

func TestInspector(t *testing.T) {
	app := app.New()
	app.TPS = 300

	posID := ecs.ComponentID[Position](&app.World)
	velID := ecs.ComponentID[Velocity](&app.World)
	entity := app.World.Unsafe().NewEntity(posID, velID)

	ecs.AddResource(&app.World, &resource.SelectedEntity{Selected: entity})

	app.AddUISystem((&window.Window{}).
		With(&monitor.Inspector{
			HideNames: true,
		}))

	app.AddSystem(&system.FixedTermination{
		Steps: 100,
	})

	app.Run()
}

func TestInspector_DeadEntity(t *testing.T) {
	app := app.New()
	app.TPS = 300

	posID := ecs.ComponentID[Position](&app.World)
	velID := ecs.ComponentID[Velocity](&app.World)
	entity := app.World.Unsafe().NewEntity(posID, velID)

	ecs.AddResource(&app.World, &resource.SelectedEntity{Selected: entity})

	app.AddUISystem((&window.Window{}).
		With(&monitor.Inspector{
			HideNames: true,
		}))

	app.AddSystem(&system.FixedTermination{
		Steps: 100,
	})

	app.World.RemoveEntity(entity)

	app.Run()
}
