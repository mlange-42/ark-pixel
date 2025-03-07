package monitor

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

func TestRingBuffer(t *testing.T) {
	r := newRingBuffer[int](10)

	for i := 0; i < 10; i++ {
		r.Add(i)
		assert.Equal(t, i+1, r.Len())
	}

	for i := 0; i < 10; i++ {
		assert.Equal(t, i, r.Get(i))
	}

	for i := 0; i < 10; i++ {
		r.Add(i + 10)
		assert.Equal(t, 10, r.Len())
		assert.Equal(t, i+1, r.Get(0))
	}

	for i := 0; i < 10; i++ {
		assert.Equal(t, i+10, r.Get(i))
	}
}

func TestCalcTps(t *testing.T) {
	tps := calcTps(1, true)
	assert.Equal(t, 2.0, tps)

	tps = calcTps(1, false)
	assert.Equal(t, 0.0, tps)

	tps = calcTps(100, true)
	assert.Equal(t, 120.0, tps)

	tps = calcTps(100, false)
	assert.Equal(t, 80.0, tps)

	tps = calcTps(99, true)
	assert.Equal(t, 100.0, tps)

	tps = calcTps(99, false)
	assert.Equal(t, 80.0, tps)

	tps = calcTps(12345, true)
	assert.Equal(t, 12345.0, tps)
}
