package log

import (
	"fmt"
	"os"
	"sync/atomic"
	"time"

	"github.com/derision-test/glock"
)

type logSink interface {
	Log(timestamp time.Time, level LogLevel, fields LogFields, msg string) error
}

type baseWrapper struct {
	logSink  logSink
	level    LogLevel
	clock    glock.Clock
	exiter   func()
	sequence uint64
}

type baseLogger struct {
	wrapper *baseWrapper
	fields  LogFields
}

func newBaseLogger(logSink logSink, level LogLevel, initialFields LogFields) Logger {
	wrapper := &baseWrapper{
		logSink,
		level,
		glock.NewRealClock(),
		func() { os.Exit(1) },
		0,
	}

	return FromMinimalLogger(&baseLogger{wrapper, initialFields})
}

func newTestLogger(logSink logSink, level LogLevel, initialFields LogFields, clock glock.Clock, exiter func()) Logger {
	wrapper := &baseWrapper{
		logSink,
		level,
		clock,
		exiter,
		0,
	}

	return FromMinimalLogger(&baseLogger{wrapper, initialFields})
}

func (s *baseLogger) WithFields(fields LogFields) MinimalLogger {
	if len(fields) == 0 {
		return s
	}

	return &baseLogger{s.wrapper, s.fields.concat(fields)}
}

func (s *baseLogger) LogWithFields(level LogLevel, fields LogFields, format string, args ...interface{}) {
	if level > s.wrapper.level {
		return
	}

	seq := atomic.AddUint64(&s.wrapper.sequence, 1)
	fields = fields.normalizeTimeValues()
	fields["sequenceNumber"] = seq

	s.wrapper.logSink.Log(
		s.wrapper.clock.Now().UTC(),
		level,
		s.fields.concat(fields),
		fmt.Sprintf(format, args...),
	)

	if level == LevelFatal {
		s.wrapper.exiter()
	}
}

func (s *baseLogger) Sync() error {
	return nil
}
