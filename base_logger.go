package log

import (
	"fmt"
	"os"
	"sync/atomic"
	"time"

	"github.com/derision-test/glock"
)

type baseLogger interface {
	Log(timestamp time.Time, level LogLevel, fields LogFields, msg string) error
}

type baseWrapper struct {
	logger   baseLogger
	level    LogLevel
	clock    glock.Clock
	exiter   func()
	sequence uint64
}

type baseShim struct {
	wrapper *baseWrapper
	fields  LogFields
}

//
// Shim

func newBaseShim(logger baseLogger, level LogLevel, initialFields LogFields) Logger {
	wrapper := &baseWrapper{
		logger,
		level,
		glock.NewRealClock(),
		func() { os.Exit(1) },
		0,
	}

	return adaptShim(&baseShim{wrapper, initialFields})
}

func newTestShim(logger baseLogger, level LogLevel, initialFields LogFields, clock glock.Clock, exiter func()) Logger {
	wrapper := &baseWrapper{
		logger,
		level,
		clock,
		exiter,
		0,
	}

	return adaptShim(&baseShim{wrapper, initialFields})
}

func (s *baseShim) WithFields(fields LogFields) logShim {
	if len(fields) == 0 {
		return s
	}

	return &baseShim{s.wrapper, s.fields.concat(fields)}
}

func (s *baseShim) LogWithFields(level LogLevel, fields LogFields, format string, args ...interface{}) {
	if level > s.wrapper.level {
		return
	}

	seq := atomic.AddUint64(&s.wrapper.sequence, 1)
	fields = fields.normalizeTimeValues()
	fields["sequenceNumber"] = seq

	s.wrapper.logger.Log(
		s.wrapper.clock.Now().UTC(),
		level,
		s.fields.concat(fields),
		fmt.Sprintf(format, args...),
	)

	if level == LevelFatal {
		s.wrapper.exiter()
	}
}

func (s *baseShim) Sync() error {
	return nil
}
