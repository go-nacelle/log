package log

type (
	logShim interface {
		WithFields(LogFields) logShim
		LogWithFields(LogLevel, LogFields, string, ...interface{})
		Sync() error
	}

	shimAdapter struct {
		shim  logShim
		depth int
	}

	replayShimAdapter struct {
		Logger
		shim *replayShim
	}

	logMessage struct {
		level  LogLevel
		fields LogFields
		format string
		args   []interface{}
	}
)

func adaptShim(shim logShim) Logger {
	return &shimAdapter{shim: shim}
}

func adaptReplayShim(shim *replayShim) ReplayLogger {
	return &replayShimAdapter{adaptShim(shim), shim}
}

func (sa *shimAdapter) WithIndirectCaller(frames int) Logger {
	if frames <= 0 {
		panic("WithIndirectCaller called with invalid frame count")
	}

	return &shimAdapter{shim: sa.shim, depth: sa.depth + frames}
}

func (sa *shimAdapter) WithFields(fields LogFields) Logger {
	if len(fields) == 0 {
		return sa
	}

	return &shimAdapter{shim: sa.shim.WithFields(fields)}
}

func (sa *shimAdapter) LogWithFields(level LogLevel, fields LogFields, format string, args ...interface{}) {
	sa.shim.LogWithFields(level, addCaller(fields, sa.depth), format, args...)
}

func (sa *shimAdapter) Sync() error {
	return sa.shim.Sync()
}

func (sa *shimAdapter) Debug(format string, args ...interface{}) {
	sa.shim.LogWithFields(LevelDebug, addCaller(nil, sa.depth), format, args...)
}

func (sa *shimAdapter) Info(format string, args ...interface{}) {
	sa.shim.LogWithFields(LevelInfo, addCaller(nil, sa.depth), format, args...)
}

func (sa *shimAdapter) Warning(format string, args ...interface{}) {
	sa.shim.LogWithFields(LevelWarning, addCaller(nil, sa.depth), format, args...)
}

func (sa *shimAdapter) Error(format string, args ...interface{}) {
	sa.shim.LogWithFields(LevelError, addCaller(nil, sa.depth), format, args...)
}

func (sa *shimAdapter) Fatal(format string, args ...interface{}) {
	sa.shim.LogWithFields(LevelFatal, addCaller(nil, sa.depth), format, args...)
}

func (sa *shimAdapter) DebugWithFields(fields LogFields, format string, args ...interface{}) {
	sa.shim.LogWithFields(LevelDebug, addCaller(fields, sa.depth), format, args...)
}

func (sa *shimAdapter) InfoWithFields(fields LogFields, format string, args ...interface{}) {
	sa.shim.LogWithFields(LevelInfo, addCaller(fields, sa.depth), format, args...)
}

func (sa *shimAdapter) WarningWithFields(fields LogFields, format string, args ...interface{}) {
	sa.shim.LogWithFields(LevelWarning, addCaller(fields, sa.depth), format, args...)
}

func (sa *shimAdapter) ErrorWithFields(fields LogFields, format string, args ...interface{}) {
	sa.shim.LogWithFields(LevelError, addCaller(fields, sa.depth), format, args...)
}

func (sa *shimAdapter) FatalWithFields(fields LogFields, format string, args ...interface{}) {
	sa.shim.LogWithFields(LevelFatal, addCaller(fields, sa.depth), format, args...)
}

func (a *replayShimAdapter) Replay(level LogLevel) {
	a.shim.Replay(level)
}
