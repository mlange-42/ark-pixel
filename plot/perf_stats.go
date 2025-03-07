package plot

import (
	"fmt"
	"image/color"
	"time"

	px "github.com/gopxl/pixel/v2"
	"github.com/gopxl/pixel/v2/backends/opengl"
	"github.com/gopxl/pixel/v2/ext/imdraw"
	"github.com/gopxl/pixel/v2/ext/text"
	"github.com/mlange-42/ark/ecs"
)

// PerfStats drawer for performance statistics.
//
// Adds an overlay with performance statistics in the top left corner of the window.
type PerfStats struct {
	SampleInterval time.Duration // Approx. time between measurements. Optional, default 1 second.
	drawer         imdraw.IMDraw
	stats          tempStats
	summary        *text.Text
	frameTimer     frameTimer
	startTime      time.Time
	lastPlotUpdate time.Time
	step           int64
}

// Initialize the system
func (p *PerfStats) Initialize(w *ecs.World, win *opengl.Window) {
	if p.SampleInterval <= 0 {
		p.SampleInterval = time.Second
	}
	p.lastPlotUpdate = time.Now()
	p.startTime = p.lastPlotUpdate

	p.drawer = *imdraw.New(nil)

	p.summary = text.New(px.V(0, 0), defaultFont)
	p.summary.AlignedTo(px.BottomRight)

	p.step = 0

	st := w.Stats()
	p.stats.Entities = st.Entities.Used
	p.stats.Mem = st.Memory
}

// Update the drawer.
func (p *PerfStats) Update(w *ecs.World) {
	t := time.Now()
	p.frameTimer.Update(p.step, t)

	if t.Sub(p.lastPlotUpdate) >= p.SampleInterval {
		st := w.Stats()
		p.stats.Entities = st.Entities.Used
		p.stats.Mem = st.Memory
		p.lastPlotUpdate = t
	}

	p.step++
}

// UpdateInputs handles input events of the previous frame update.
func (p *PerfStats) UpdateInputs(w *ecs.World, win *opengl.Window) {}

// Draw the system
func (p *PerfStats) Draw(w *ecs.World, win *opengl.Window) {
	p.summary.Clear()
	mem, units := toMemText(p.stats.Mem)
	fmt.Fprintf(
		p.summary, "Tick: %7d\nEnt.: %7d\nTPS: %8.1f\nTPT: %6.2fms\nMem: %6.1f%s\nTime: %7s",
		p.step, p.stats.Entities, p.frameTimer.FPS(),
		float64(p.frameTimer.FrameTime().Microseconds())/1000,
		mem, units, time.Since(p.startTime).Round(time.Second),
	)

	dr := &p.drawer
	height := win.Canvas().Bounds().H()
	x0 := 10.0
	y0 := height - 10.0

	v1 := px.V(x0+p.summary.Bounds().Min.X-5, y0+p.summary.Bounds().Min.Y-12)
	v2 := px.V(x0+p.summary.Bounds().Max.X+5, y0+p.summary.Bounds().Max.Y-8)

	dr.Color = color.Black
	dr.Push(v1, v2)
	dr.Rectangle(0)

	dr.Color = color.White
	dr.Push(v1, v2)
	dr.Rectangle(1)

	dr.Draw(win)
	dr.Reset()
	dr.Clear()

	p.summary.Draw(win, px.IM.Moved(px.V(x0, y0)))

}

type tempStats struct {
	Entities int
	Mem      int
}
