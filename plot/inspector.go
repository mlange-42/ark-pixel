package plot

import (
	"fmt"
	"io"
	"reflect"

	px "github.com/gopxl/pixel/v2"
	"github.com/gopxl/pixel/v2/backends/opengl"
	"github.com/gopxl/pixel/v2/ext/text"
	"github.com/mlange-42/ark-tools/resource"
	"github.com/mlange-42/ark/ecs"
)

// Inspector drawer for inspecting entities.
//
// Shows information of the entity indicated by the SelectedEntity resource ([github.com/mlange-42/arche-model/resource.SelectedEntity]).
// Entity selection is to be done by another system, e.g. by user input.
//
// Details can be adjusted using the HideXxx fields.
// Further, keys F, T, V and N can be used to toggle details during a running simulation.
// The view can be scrolled using arrow keys or the mouse wheel.
type Inspector struct {
	HideFields  bool // Hides components fields.
	HideTypes   bool // Hides field types.
	HideValues  bool // Hides field values.
	HideNames   bool // Hide field names of nested structs.
	scroll      int
	selectedRes ecs.Resource[resource.SelectedEntity]
	text        *text.Text
	helpText    *text.Text
}

// Initialize the system
func (i *Inspector) Initialize(w *ecs.World, win *opengl.Window) {
	i.selectedRes = ecs.NewResource[resource.SelectedEntity](w)

	i.text = text.New(px.V(0, 0), defaultFont)
	i.helpText = text.New(px.V(0, 0), defaultFont)

	i.text.AlignedTo(px.BottomRight)
	i.helpText.AlignedTo(px.BottomRight)

	fmt.Fprint(i.helpText, "Toggle [f]ields, [t]ypes, [v]alues or [n]ames, scroll with arrows or mouse wheel.")
}

// Update the drawer.
func (i *Inspector) Update(w *ecs.World) {}

// UpdateInputs handles input events of the previous frame update.
func (i *Inspector) UpdateInputs(w *ecs.World, win *opengl.Window) {
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
func (i *Inspector) Draw(w *ecs.World, win *opengl.Window) {
	i.helpText.Draw(win, px.IM.Moved(px.V(10, 20)))

	if !i.selectedRes.Has() {
		return
	}
	sel := i.selectedRes.Get().Selected
	if sel.IsZero() {
		return
	}

	height := win.Canvas().Bounds().H()
	x0 := 10.0
	y0 := height - 10.0

	i.text.Clear()
	fmt.Fprintf(i.text, "Entity %+v\n\n", sel)

	if !w.Alive(sel) {
		fmt.Fprint(i.text, "  dead entity")
		i.text.Draw(win, px.IM.Moved(px.V(x0, y0)))
		return
	}

	scroll := i.scroll
	ids := w.Unsafe().IDs(sel)
	for _, id := range ids {
		tp, _ := ecs.ComponentInfo(w, id)
		ptr := w.Unsafe().Get(sel, id)
		val := reflect.NewAt(tp.Type, ptr).Elem()

		if scroll <= 0 {
			fmt.Fprintf(i.text, "  %s\n", tp.Type.Name())
		}
		scroll--

		if !i.HideFields {
			for k := 0; k < val.NumField(); k++ {
				field := tp.Type.Field(k)
				if field.IsExported() {
					if scroll <= 0 {
						i.printField(i.text, tp.Type, field, val.Field(k))
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

func (i *Inspector) printField(w io.Writer, tp reflect.Type, field reflect.StructField, value reflect.Value) {
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
