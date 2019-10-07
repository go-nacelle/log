package log

//go:generate go-mockgen -f github.com/go-nacelle/log -i baseLogger -o base_logger_mock_test.go

import (
	"github.com/aphistic/sweet"
	"github.com/efritz/glock"
	. "github.com/efritz/go-mockgen/matchers"
	. "github.com/onsi/gomega"
)

type BaseLoggerSuite struct{}

func (s *BaseLoggerSuite) TestLogFormat(t sweet.T) {
	base := NewMockBaseLogger()
	clock := glock.NewMockClock()
	logger := newTestShim(base, LevelDebug, nil, clock, func() {})
	logger.LogWithFields(LevelInfo, nil, "test %d %d %d", 1, 2, 3)

	Expect(base.LogFunc).To(BeCalledOnceWith(
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

func (s *BaseLoggerSuite) TestWrappedLoggers(t sweet.T) {
	base := NewMockBaseLogger()
	clock := glock.NewMockClock()
	logger := newTestShim(base, LevelDebug, LogFields{"init": "foo"}, clock, func() {})
	wrappedLogger := logger.WithFields(LogFields{"wrapped": "bar"})
	wrappedLogger.LogWithFields(LevelDebug, LogFields{"extra": "baz"}, "test %d %d %d", 1, 2, 3)
	logger.LogWithFields(LevelDebug, LogFields{"extra": "bonk"}, "test %d %d %d", 1, 2, 3)

	Expect(base.LogFunc).To(BeCalledOnceWith(
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

	Expect(base.LogFunc).To(BeCalledOnceWith(
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

func (s *BaseLoggerSuite) TestLogLevelFilter(t sweet.T) {
	base := NewMockBaseLogger()
	clock := glock.NewMockClock()
	logger := newTestShim(base, LevelInfo, nil, clock, func() {})
	logger.LogWithFields(LevelDebug, nil, "test %d %d %d", 1, 2, 3)
	Expect(base.LogFunc).NotTo(BeCalled())
}

func (s *BaseLoggerSuite) TestLogFatal(t sweet.T) {
	base := NewMockBaseLogger()
	clock := glock.NewMockClock()
	called := false
	logger := newTestShim(base, LevelInfo, nil, clock, func() { called = true })
	logger.LogWithFields(LevelFatal, nil, "test %d %d %d", 1, 2, 3)
	Expect(called).To(BeTrue())
}
