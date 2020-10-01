package log

import (
	"testing"
	"time"

	"github.com/derision-test/glock"
	"github.com/stretchr/testify/assert"
)

func TestRollupAdapterSimilarMessages(t *testing.T) {
	shim := &testShim{}
	clock := glock.NewMockClock()
	adapter := newRollupShim(adaptShim(shim), clock, time.Second)

	for i := 1; i <= 20; i++ {
		// Logged, starting window
		adapter.LogWithFields(LevelDebug, nil, "a")
		assert.Len(t, shim.copy(), 2*i-1)

		// Stashed
		adapter.LogWithFields(LevelDebug, nil, "a")
		adapter.LogWithFields(LevelDebug, nil, "a")
		assert.Len(t, shim.copy(), 2*i-1)

		// Flushed
		clock.BlockingAdvance(time.Second)
		requireEventually(t, func() bool { return len(shim.copy()) == 2*i })
		assert.Equal(t, 2, shim.copy()[2*i-1].fields[FieldRollup])
	}
}

func TestRollupAdapterInactivity(t *testing.T) {
	shim := &testShim{}
	clock := glock.NewMockClock()
	adapter := newRollupShim(adaptShim(shim), clock, time.Second)

	for i := 0; i < 20; i++ {
		adapter.LogWithFields(LevelDebug, nil, "a")
		clock.Advance(time.Second * 2)
	}

	// All messages present
	eventually(t, func() bool { return len(shim.copy()) == 20 })
}

func TestRollupAdapterFlushesRelativeToFirstMessage(t *testing.T) {
	shim := &testShim{}
	clock := glock.NewMockClock()
	adapter := newRollupShim(adaptShim(shim), clock, time.Second)

	adapter.LogWithFields(LevelDebug, nil, "a")
	clock.Advance(time.Millisecond * 500)

	for i := 0; i < 90; i++ {
		adapter.LogWithFields(LevelDebug, nil, "a")
		clock.Advance(time.Millisecond * 5)
	}

	clock.BlockingAdvance(time.Millisecond * 50)
	eventually(t, func() bool { return len(shim.copy()) == 2 })
}

func TestRollupAdapterAllDistinctMessages(t *testing.T) {
	shim := &testShim{}
	clock := glock.NewMockClock()
	adapter := newRollupShim(adaptShim(shim), clock, time.Second)

	for i := 0; i < 10; i++ {
		adapter.LogWithFields(LevelDebug, nil, "a")
		adapter.LogWithFields(LevelDebug, nil, "b")
		adapter.LogWithFields(LevelDebug, nil, "c")
	}

	assert.Len(t, shim.copy(), 3)
}
