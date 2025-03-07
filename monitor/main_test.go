package monitor_test

import (
	"os"
	"testing"

	"github.com/gopxl/pixel/v2/backends/opengl"
)

func TestMain(m *testing.M) {
	opengl.Run(func() {
		os.Exit(m.Run())
	})
}
