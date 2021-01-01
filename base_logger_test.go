package log

import (
	"testing"

	"github.com/derision-test/glock"
	mockassert "github.com/derision-test/go-mockgen/testutil/assert"
	"github.com/stretchr/testify/assert"
)

type BaseLoggerSuite struct{}

func (s *BaseLoggerSuite) TestLogFormat(t *testing.T) {
	base := NewMockBaseLogger()
	clock := glock.NewMockClock()
	logger := newTestLogger(base, LevelDebug, nil, clock, func() {})
	logger.LogWithFields(LevelInfo, nil, "test %d %d %d", 1, 2, 3)

	mockassert.CalledOnceWith(t, base.LogFunc, mockassert.Values(
		clock.Now().UTC(),
		LevelInfo,
		LogFields{
			// Note: this value refers to the line number containing `LogWithFields` in
			// the test setup above. If code is added before that line, this value must
			// be updated.
			"caller":         "log/base_logger_test.go:18",
			"sequenceNumber": uint64(1),
		},
		"test 1 2 3",
	))
}

func (s *BaseLoggerSuite) TestWrappedLoggers(t *testing.T) {
	base := NewMockBaseLogger()
	clock := glock.NewMockClock()
	logger := newTestLogger(base, LevelDebug, LogFields{"init": "foo"}, clock, func() {})
	wrappedLogger := logger.WithFields(LogFields{"wrapped": "bar"})
	wrappedLogger.LogWithFields(LevelDebug, LogFields{"extra": "baz"}, "test %d %d %d", 1, 2, 3)
	logger.LogWithFields(LevelDebug, LogFields{"extra": "bonk"}, "test %d %d %d", 1, 2, 3)

	mockassert.CalledOnceWith(t, base.LogFunc, mockassert.Values(
		clock.Now().UTC(),
		LevelDebug,
		LogFields{
			"init":    "foo",
			"wrapped": "bar",
			"extra":   "baz",
			// Note: this value refers to the line number containing `LogWithFields` in
			// the test setup above. If code is added before that line, this value must
			// be updated.
			"caller":         "log/base_logger_test.go:39",
			"sequenceNumber": uint64(1),
		},
		"test 1 2 3",
	))

	mockassert.CalledOnceWith(t, base.LogFunc, mockassert.Values(
		clock.Now().UTC(),
		LevelDebug,
		LogFields{
			"init":  "foo",
			"extra": "bonk",
			// Note: this value refers to the line number containing `LogWithFields` in
			// the test setup above. If code is added before that line, this value must
			// be updated.
			"caller":         "log/base_logger_test.go:40",
			"sequenceNumber": uint64(2),
		},
		"test 1 2 3",
	))
}

func (s *BaseLoggerSuite) TestLogLevelFilter(t *testing.T) {
	base := NewMockBaseLogger()
	clock := glock.NewMockClock()
	logger := newTestLogger(base, LevelInfo, nil, clock, func() {})
	logger.LogWithFields(LevelDebug, nil, "test %d %d %d", 1, 2, 3)
	mockassert.NotCalled(t, base.LogFunc)
}

func (s *BaseLoggerSuite) TestLogFatal(t *testing.T) {
	base := NewMockBaseLogger()
	clock := glock.NewMockClock()
	called := false
	logger := newTestLogger(base, LevelInfo, nil, clock, func() { called = true })
	logger.LogWithFields(LevelFatal, nil, "test %d %d %d", 1, 2, 3)
	assert.True(t, called)
}
