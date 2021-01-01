package log

import (
	"testing"

	"github.com/derision-test/glock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReplayLogger(t *testing.T) {
	loggerm := &testLogger{}
	clock := glock.NewMockClock()
	replayLogger := newReplayLogger(FromMinimalLogger(loggerm), clock, LevelDebug)

	replayLogger.LogWithFields(LevelDebug, LogFields{"x": "x"}, "foo", 12)
	replayLogger.LogWithFields(LevelDebug, LogFields{"y": "y"}, "bar", 43)
	replayLogger.LogWithFields(LevelDebug, LogFields{"z": "z"}, "baz", 74)
	replayLogger.Replay(LevelWarning)

	messages := loggerm.copy()
	require.Len(t, messages, 6)

	for i := 0; i < 3; i++ {
		assert.Equal(t, LevelDebug, messages[i+0].level)
		assert.Equal(t, LevelWarning, messages[i+3].level)
	}

	for i, format := range []string{"foo", "bar", "baz"} {
		assert.Equal(t, format, messages[i+0].format)
		assert.Equal(t, format, messages[i+3].format)
	}

	for i, expected := range []int{12, 43, 74} {
		assert.Equal(t, expected, messages[i+0].args[0])
		assert.Equal(t, expected, messages[i+3].args[0])
	}

	for i, field := range []string{"x", "y", "z"} {
		assert.Equal(t, field, messages[i+0].fields[field])
		assert.Equal(t, field, messages[i+3].fields[field])
	}
}

func TestReplayLoggerTwice(t *testing.T) {
	logger := &testLogger{}
	clock := glock.NewMockClock()
	replayLogger := newReplayLogger(FromMinimalLogger(logger), clock, LevelDebug)

	replayLogger.LogWithFields(LevelDebug, nil, "foo")
	replayLogger.LogWithFields(LevelDebug, nil, "bar")
	replayLogger.LogWithFields(LevelDebug, nil, "baz")
	replayLogger.Replay(LevelWarning)
	replayLogger.Replay(LevelError)

	messages := logger.copy()
	require.Len(t, messages, 9)
	assert.Equal(t, LevelDebug, messages[0].level)
	assert.Equal(t, LevelDebug, messages[1].level)
	assert.Equal(t, LevelDebug, messages[2].level)
	assert.Equal(t, LevelWarning, messages[3].level)
	assert.Equal(t, LevelWarning, messages[4].level)
	assert.Equal(t, LevelWarning, messages[5].level)
	assert.Equal(t, LevelError, messages[6].level)
	assert.Equal(t, LevelError, messages[7].level)
	assert.Equal(t, LevelError, messages[8].level)

	for i, format := range []string{"foo", "bar", "baz", "foo", "bar", "baz", "foo", "bar", "baz"} {
		assert.Equal(t, format, messages[i].format)
	}
}

func TestReplayLoggerAtHigherlevelNoops(t *testing.T) {
	logger := &testLogger{}
	clock := glock.NewMockClock()
	replayLogger := newReplayLogger(FromMinimalLogger(logger), clock, LevelDebug)

	replayLogger.LogWithFields(LevelDebug, nil, "foo")
	replayLogger.LogWithFields(LevelDebug, nil, "bar")
	replayLogger.LogWithFields(LevelDebug, nil, "baz")
	replayLogger.Replay(LevelError)
	replayLogger.Replay(LevelWarning)

	messages := logger.copy()
	require.Len(t, messages, 6)
	assert.Equal(t, LevelDebug, messages[0].level)
	assert.Equal(t, LevelDebug, messages[1].level)
	assert.Equal(t, LevelDebug, messages[2].level)
	assert.Equal(t, LevelError, messages[3].level)
	assert.Equal(t, LevelError, messages[4].level)
	assert.Equal(t, LevelError, messages[5].level)

	for i, format := range []string{"foo", "bar", "baz", "foo", "bar", "baz"} {
		assert.Equal(t, format, messages[i].format)
	}
}

func TestReplayLoggerLogAfterReplaySendsImmediately(t *testing.T) {
	logger := &testLogger{}
	clock := glock.NewMockClock()
	replayLogger := newReplayLogger(FromMinimalLogger(logger), clock, LevelDebug)

	replayLogger.LogWithFields(LevelDebug, nil, "foo")
	replayLogger.LogWithFields(LevelDebug, nil, "bar")
	replayLogger.LogWithFields(LevelDebug, nil, "baz")
	replayLogger.Replay(LevelWarning)
	replayLogger.LogWithFields(LevelDebug, nil, "bnk")
	replayLogger.LogWithFields(LevelDebug, nil, "qux")

	messages := logger.copy()
	require.Len(t, messages, 10)
	assert.Equal(t, LevelDebug, messages[0].level)
	assert.Equal(t, LevelDebug, messages[1].level)
	assert.Equal(t, LevelDebug, messages[2].level)
	assert.Equal(t, LevelWarning, messages[3].level)
	assert.Equal(t, LevelWarning, messages[4].level)
	assert.Equal(t, LevelWarning, messages[5].level)
	assert.Equal(t, LevelDebug, messages[6].level)
	assert.Equal(t, LevelWarning, messages[7].level)
	assert.Equal(t, LevelDebug, messages[8].level)
	assert.Equal(t, LevelWarning, messages[9].level)

	for i, format := range []string{"foo", "bar", "baz", "foo", "bar", "baz", "bnk", "bnk", "qux", "qux"} {
		assert.Equal(t, format, messages[i].format)
	}
}

func TestReplayLoggerLogAfterSecondReplaySendsAtNewLevel(t *testing.T) {
	logger := &testLogger{}
	clock := glock.NewMockClock()
	replayLogger := newReplayLogger(FromMinimalLogger(logger), clock, LevelDebug)

	replayLogger.LogWithFields(LevelDebug, nil, "foo")
	replayLogger.LogWithFields(LevelDebug, nil, "bar")
	replayLogger.Replay(LevelWarning)
	replayLogger.Replay(LevelError)
	replayLogger.LogWithFields(LevelDebug, nil, "baz")
	replayLogger.LogWithFields(LevelDebug, nil, "bnk")

	messages := logger.copy()
	require.Len(t, messages, 10)
	assert.Equal(t, LevelDebug, messages[0].level)
	assert.Equal(t, LevelDebug, messages[1].level)
	assert.Equal(t, LevelWarning, messages[2].level)
	assert.Equal(t, LevelWarning, messages[3].level)
	assert.Equal(t, LevelError, messages[4].level)
	assert.Equal(t, LevelError, messages[5].level)
	assert.Equal(t, LevelDebug, messages[6].level)
	assert.Equal(t, LevelError, messages[7].level)
	assert.Equal(t, LevelDebug, messages[8].level)
	assert.Equal(t, LevelError, messages[9].level)

	for i, format := range []string{"foo", "bar", "foo", "bar", "foo", "bar", "baz", "baz", "bnk", "bnk"} {
		assert.Equal(t, format, messages[i].format)
	}
}

func TestReplayLoggerCheckReplayAddsAttribute(t *testing.T) {
	logger := &testLogger{}
	clock := glock.NewMockClock()
	replayLogger := newReplayLogger(FromMinimalLogger(logger), clock, LevelDebug, LevelInfo)

	replayLogger.LogWithFields(LevelDebug, nil, "foo")
	replayLogger.LogWithFields(LevelInfo, nil, "bar")
	replayLogger.LogWithFields(LevelDebug, nil, "baz")
	replayLogger.Replay(LevelError)
	replayLogger.LogWithFields(LevelDebug, nil, "bonk")

	messages := logger.copy()
	require.Len(t, messages, 8)
	assert.NotContains(t, messages[0].fields, FieldReplay)
	assert.NotContains(t, messages[1].fields, FieldReplay)
	assert.NotContains(t, messages[2].fields, FieldReplay)
	assert.Equal(t, LevelDebug, messages[3].fields[FieldReplay])
	assert.Equal(t, LevelInfo, messages[4].fields[FieldReplay])
	assert.Equal(t, LevelDebug, messages[5].fields[FieldReplay])
	assert.NotContains(t, messages[6].fields, FieldReplay)
	assert.Equal(t, LevelDebug, messages[7].fields[FieldReplay])
}

func TestReplayLoggerCheckSecondReplayAddsAttribute(t *testing.T) {
	logger := &testLogger{}
	clock := glock.NewMockClock()
	replayLogger := newReplayLogger(FromMinimalLogger(logger), clock, LevelDebug, LevelInfo)

	replayLogger.LogWithFields(LevelDebug, nil, "foo")
	replayLogger.LogWithFields(LevelInfo, nil, "bar")
	replayLogger.Replay(LevelWarning)
	replayLogger.Replay(LevelError)
	replayLogger.LogWithFields(LevelDebug, nil, "bnk")

	messages := logger.copy()
	require.Len(t, messages, 8)
	assert.NotContains(t, messages[0].fields, FieldReplay)
	assert.NotContains(t, messages[1].fields, FieldReplay)
	assert.Equal(t, LevelDebug, messages[2].fields[FieldReplay])
	assert.Equal(t, LevelInfo, messages[3].fields[FieldReplay])
	assert.Equal(t, LevelDebug, messages[4].fields[FieldReplay])
	assert.Equal(t, LevelInfo, messages[5].fields[FieldReplay])
	assert.NotContains(t, messages[6].fields, FieldReplay)
	assert.Equal(t, LevelDebug, messages[7].fields[FieldReplay])
}
