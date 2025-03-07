package plot

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFind(t *testing.T) {
	arr := []string{"A", "B", "C"}

	idx, ok := find(arr, "B")
	assert.True(t, ok)
	assert.Equal(t, 1, idx)

	idx, ok = find(arr, "D")
	assert.False(t, ok)
	assert.Equal(t, -1, idx)
}
