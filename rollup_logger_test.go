package log

import (
	"testing"
	"time"

	"github.com/derision-test/glock"
	"github.com/stretchr/testify/assert"
)

func TestRollupLoggerSimilarMessages(t *testing.T) {
	logger := &testLogger{}
	clock := glock.NewMockClock()
	rollupLogger := newRollupLogger(FromMinimalLogger(logger), clock, time.Second)

	for i := 1; i <= 20; i++ {
		// Logged, starting window
		rollupLogger.LogWithFields(LevelDebug, nil, "a")
		assert.Len(t, logger.copy(), 2*i-1)

		// Stashed
		rollupLogger.LogWithFields(LevelDebug, nil, "a")
		rollupLogger.LogWithFields(LevelDebug, nil, "a")
		assert.Len(t, logger.copy(), 2*i-1)

		// Flushed
		clock.BlockingAdvance(time.Second)
		requireEventually(t, func() bool { return len(logger.copy()) == 2*i })
		assert.Equal(t, 2, logger.copy()[2*i-1].fields[FieldRollup])
	}
}

func TestRollupLoggerInactivity(t *testing.T) {
	logger := &testLogger{}
	clock := glock.NewMockClock()
	rollupLogger := newRollupLogger(FromMinimalLogger(logger), clock, time.Second)

	for i := 0; i < 20; i++ {
		rollupLogger.LogWithFields(LevelDebug, nil, "a")
		clock.Advance(time.Second * 2)
	}

	// All messages present
	eventually(t, func() bool { return len(logger.copy()) == 20 })
}

func TestRollupLoggerFlushesRelativeToFirstMessage(t *testing.T) {
	logger := &testLogger{}
	clock := glock.NewMockClock()
	rollupLogger := newRollupLogger(FromMinimalLogger(logger), clock, time.Second)

	rollupLogger.LogWithFields(LevelDebug, nil, "a")
	clock.Advance(time.Millisecond * 500)

	for i := 0; i < 90; i++ {
		rollupLogger.LogWithFields(LevelDebug, nil, "a")
		clock.Advance(time.Millisecond * 5)
	}

	clock.BlockingAdvance(time.Millisecond * 50)
	eventually(t, func() bool { return len(logger.copy()) == 2 })
}

func TestRollupLoggerAllDistinctMessages(t *testing.T) {
	logger := &testLogger{}
	clock := glock.NewMockClock()
	rollupLogger := newRollupLogger(FromMinimalLogger(logger), clock, time.Second)

	for i := 0; i < 10; i++ {
		rollupLogger.LogWithFields(LevelDebug, nil, "a")
		rollupLogger.LogWithFields(LevelDebug, nil, "b")
		rollupLogger.LogWithFields(LevelDebug, nil, "c")
	}

	assert.Len(t, logger.copy(), 3)
}
