package window_test

import (
	"testing"

	"github.com/gopxl/pixel/v2"
	"github.com/gopxl/pixel/v2/backends/opengl"
	"github.com/mlange-42/ark-pixel/window"
	"github.com/stretchr/testify/assert"
)

func TestScale(t *testing.T) {
	w, err := opengl.NewWindow(opengl.WindowConfig{Bounds: pixel.R(0, 0, 800, 600)})

	if err != nil {
		panic(err)
	}

	scale := window.Scale(w, 400, 300)

	assert.Equal(t, 2.0, scale)
}
