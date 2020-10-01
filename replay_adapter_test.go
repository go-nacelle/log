package log

import (
	"testing"

	"github.com/derision-test/glock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReplayAdapter(t *testing.T) {
	shim := &testShim{}
	clock := glock.NewMockClock()
	adapter := newReplayShim(adaptShim(shim), clock, LevelDebug)

	adapter.LogWithFields(LevelDebug, LogFields{"x": "x"}, "foo", 12)
	adapter.LogWithFields(LevelDebug, LogFields{"y": "y"}, "bar", 43)
	adapter.LogWithFields(LevelDebug, LogFields{"z": "z"}, "baz", 74)
	adapter.Replay(LevelWarning)

	messages := shim.copy()
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

func TestReplayAdapterTwice(t *testing.T) {
	shim := &testShim{}
	clock := glock.NewMockClock()
	adapter := newReplayShim(adaptShim(shim), clock, LevelDebug)

	adapter.LogWithFields(LevelDebug, nil, "foo")
	adapter.LogWithFields(LevelDebug, nil, "bar")
	adapter.LogWithFields(LevelDebug, nil, "baz")
	adapter.Replay(LevelWarning)
	adapter.Replay(LevelError)

	messages := shim.copy()
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

func TestReplayAdapterAtHigherlevelNoops(t *testing.T) {
	shim := &testShim{}
	clock := glock.NewMockClock()
	adapter := newReplayShim(adaptShim(shim), clock, LevelDebug)

	adapter.LogWithFields(LevelDebug, nil, "foo")
	adapter.LogWithFields(LevelDebug, nil, "bar")
	adapter.LogWithFields(LevelDebug, nil, "baz")
	adapter.Replay(LevelError)
	adapter.Replay(LevelWarning)

	messages := shim.copy()
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

func TestReplayAdapterLogAfterReplaySendsImmediately(t *testing.T) {
	shim := &testShim{}
	clock := glock.NewMockClock()
	adapter := newReplayShim(adaptShim(shim), clock, LevelDebug)

	adapter.LogWithFields(LevelDebug, nil, "foo")
	adapter.LogWithFields(LevelDebug, nil, "bar")
	adapter.LogWithFields(LevelDebug, nil, "baz")
	adapter.Replay(LevelWarning)
	adapter.LogWithFields(LevelDebug, nil, "bnk")
	adapter.LogWithFields(LevelDebug, nil, "qux")

	messages := shim.copy()
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

func TestReplayAdapterLogAfterSecondReplaySendsAtNewLevel(t *testing.T) {
	shim := &testShim{}
	clock := glock.NewMockClock()
	adapter := newReplayShim(adaptShim(shim), clock, LevelDebug)

	adapter.LogWithFields(LevelDebug, nil, "foo")
	adapter.LogWithFields(LevelDebug, nil, "bar")
	adapter.Replay(LevelWarning)
	adapter.Replay(LevelError)
	adapter.LogWithFields(LevelDebug, nil, "baz")
	adapter.LogWithFields(LevelDebug, nil, "bnk")

	messages := shim.copy()
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

func TestReplayAdapterCheckReplayAddsAttribute(t *testing.T) {
	shim := &testShim{}
	clock := glock.NewMockClock()
	adapter := newReplayShim(adaptShim(shim), clock, LevelDebug, LevelInfo)

	adapter.LogWithFields(LevelDebug, nil, "foo")
	adapter.LogWithFields(LevelInfo, nil, "bar")
	adapter.LogWithFields(LevelDebug, nil, "baz")
	adapter.Replay(LevelError)
	adapter.LogWithFields(LevelDebug, nil, "bonk")

	messages := shim.copy()
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

func TestReplayAdapterCheckSecondReplayAddsAttribute(t *testing.T) {
	shim := &testShim{}
	clock := glock.NewMockClock()
	adapter := newReplayShim(adaptShim(shim), clock, LevelDebug, LevelInfo)

	adapter.LogWithFields(LevelDebug, nil, "foo")
	adapter.LogWithFields(LevelInfo, nil, "bar")
	adapter.Replay(LevelWarning)
	adapter.Replay(LevelError)
	adapter.LogWithFields(LevelDebug, nil, "bnk")

	messages := shim.copy()
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
