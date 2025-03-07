package plot

import (
	"fmt"
	"io"
	"reflect"

	px "github.com/gopxl/pixel/v2"
	"github.com/gopxl/pixel/v2/backends/opengl"
	"github.com/gopxl/pixel/v2/ext/text"
	"github.com/mlange-42/ark-tools/app"
	"github.com/mlange-42/ark/ecs"
)

// Systems drawer for inspecting ECS systems.
//
// Lists all systems and UI systems in their scheduling order,
// with their public fields.
//
// Details can be adjusted using the HideXxx fields.
// Further, keys U, F, T, V and N can be used to toggle details during a running simulation.
// The view can be scrolled using arrow keys or the mouse wheel.
type Systems struct {
	HideUISystems bool // Hides UI systems.
	HideFields    bool // Hides components fields.
	HideTypes     bool // Hides field types.
	HideValues    bool // Hides field values.
	HideNames     bool // Hide field names of nested structs.
	scroll        int
	systemsRes    ecs.Resource[app.Systems]
	text          *text.Text
	helpText      *text.Text
}

// Initialize the system
func (i *Systems) Initialize(w *ecs.World, win *opengl.Window) {
	i.systemsRes = ecs.NewResource[app.Systems](w)

	i.text = text.New(px.V(0, 0), defaultFont)
	i.helpText = text.New(px.V(0, 0), defaultFont)

	i.text.AlignedTo(px.BottomRight)
	i.helpText.AlignedTo(px.BottomRight)

	fmt.Fprint(i.helpText, "Toggle [u]i systems, [f]ields, [t]ypes, [v]alues or [n]ames, scroll with arrows or mouse wheel.")
}

// Update the drawer.
func (i *Systems) Update(w *ecs.World) {}

// UpdateInputs handles input events of the previous frame update.
func (i *Systems) UpdateInputs(w *ecs.World, win *opengl.Window) {
	if win.JustPressed(px.KeyF) {
		i.HideFields = !i.HideFields
		return
	}
	if win.JustPressed(px.KeyT) {
		i.HideTypes = !i.HideTypes
		return
	}
	if win.JustPressed(px.KeyV) {
		i.HideValues = !i.HideValues
		return
	}
	if win.JustPressed(px.KeyN) {
		i.HideNames = !i.HideNames
		return
	}
	if win.JustPressed(px.KeyU) {
		i.HideUISystems = !i.HideUISystems
		return
	}
	if win.JustPressed(px.KeyDown) {
		i.scroll++
		return
	}
	if win.JustPressed(px.KeyUp) {
		if i.scroll > 0 {
			i.scroll--
		}
		return
	}
	scr := win.MouseScroll()
	if scr.Y != 0 {
		i.scroll -= int(scr.Y)
		if i.scroll < 0 {
			i.scroll = 0
		}
	}
}

// Draw the system
func (i *Systems) Draw(w *ecs.World, win *opengl.Window) {
	i.helpText.Draw(win, px.IM.Moved(px.V(10, 20)))

	if !i.systemsRes.Has() {
		return
	}
	systems := i.systemsRes.Get()

	height := win.Canvas().Bounds().H()
	x0 := 10.0
	y0 := height - 10.0

	i.text.Clear()
	fmt.Fprint(i.text, "Systems\n\n")

	scroll := i.scroll

	for _, sys := range systems.Systems() {
		if i.HideUISystems {
			if _, ok := sys.(app.UISystem); ok {
				continue
			}
		}

		val := reflect.ValueOf(sys).Elem()
		tp := val.Type()

		if scroll <= 0 {
			fmt.Fprintf(i.text, "  %s\n", tp.Name())
		}
		scroll--

		if !i.HideFields {
			for k := 0; k < val.NumField(); k++ {
				field := tp.Field(k)
				if field.IsExported() {
					if scroll <= 0 {
						i.printField(i.text, tp, field, val.Field(k))
					}
					scroll--
				}
			}
			if scroll <= 0 {
				fmt.Fprint(i.text, "\n")
			}
			scroll--
		}
	}

	if i.HideUISystems {
		i.text.Draw(win, px.IM.Moved(px.V(x0, y0)))
		return
	}

	fmt.Fprint(i.text, "\nUI Systems\n\n")
	for _, sys := range systems.UISystems() {
		val := reflect.ValueOf(sys).Elem()
		tp := val.Type()

		if scroll <= 0 {
			fmt.Fprintf(i.text, "  %s\n", tp.Name())
		}
		scroll--

		if !i.HideFields {
			for k := 0; k < val.NumField(); k++ {
				field := tp.Field(k)
				if field.IsExported() {
					if scroll <= 0 {
						i.printField(i.text, tp, field, val.Field(k))
					}
					scroll--
				}
			}
			if scroll <= 0 {
				fmt.Fprint(i.text, "\n")
			}
			scroll--
		}
	}

	i.text.Draw(win, px.IM.Moved(px.V(x0, y0)))
}

func (i *Systems) printField(w io.Writer, tp reflect.Type, field reflect.StructField, value reflect.Value) {
	fmt.Fprintf(w, "    %-20s ", field.Name)
	if !i.HideTypes {
		fmt.Fprintf(w, "    %-16s ", value.Type())
	}
	if !i.HideValues {
		if i.HideNames {
			fmt.Fprintf(w, "= %v", value.Interface())
		} else {
			fmt.Fprintf(w, "= %+v", value.Interface())
		}
	}
	fmt.Fprint(i.text, "\n")
}
