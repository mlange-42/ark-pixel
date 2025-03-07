package window

import (
	"github.com/gopxl/pixel/v2/backends/opengl"
	"github.com/mlange-42/ark-tools/app"
)

// Run is essential to run simulations that feature ark-pixel UI components.
// Call this function from the main function of your application, with a Model as argument.
// This is necessary, so that App.Run runs on the main thread.
func Run(app *app.App) {
	opengl.Run(app.Run)
}
