package log

type (
	MinimalLogger interface {
		WithFields(LogFields) MinimalLogger
		LogWithFields(LogLevel, LogFields, string, ...interface{})
		Sync() error
	}

	adapter struct {
		logger MinimalLogger
		depth  int
	}

	logMessage struct {
		level  LogLevel
		fields LogFields
		format string
		args   []interface{}
	}
)

func FromMinimalLogger(logger MinimalLogger) Logger {
	return &adapter{logger: logger}
}

func (sa *adapter) WithIndirectCaller(frames int) Logger {
	if frames <= 0 {
		panic("WithIndirectCaller called with invalid frame count")
	}

	return &adapter{logger: sa.logger, depth: sa.depth + frames}
}

func (sa *adapter) WithFields(fields LogFields) Logger {
	if len(fields) == 0 {
		return sa
	}

	return &adapter{logger: sa.logger.WithFields(fields)}
}

func (sa *adapter) LogWithFields(level LogLevel, fields LogFields, format string, args ...interface{}) {
	sa.logger.LogWithFields(level, addCaller(fields, sa.depth), format, args...)
}

func (sa *adapter) Sync() error {
	return sa.logger.Sync()
}

func (sa *adapter) Debug(format string, args ...interface{}) {
	sa.logger.LogWithFields(LevelDebug, addCaller(nil, sa.depth), format, args...)
}

func (sa *adapter) Info(format string, args ...interface{}) {
	sa.logger.LogWithFields(LevelInfo, addCaller(nil, sa.depth), format, args...)
}

func (sa *adapter) Warning(format string, args ...interface{}) {
	sa.logger.LogWithFields(LevelWarning, addCaller(nil, sa.depth), format, args...)
}

func (sa *adapter) Error(format string, args ...interface{}) {
	sa.logger.LogWithFields(LevelError, addCaller(nil, sa.depth), format, args...)
}

func (sa *adapter) Fatal(format string, args ...interface{}) {
	sa.logger.LogWithFields(LevelFatal, addCaller(nil, sa.depth), format, args...)
}

func (sa *adapter) DebugWithFields(fields LogFields, format string, args ...interface{}) {
	sa.logger.LogWithFields(LevelDebug, addCaller(fields, sa.depth), format, args...)
}

func (sa *adapter) InfoWithFields(fields LogFields, format string, args ...interface{}) {
	sa.logger.LogWithFields(LevelInfo, addCaller(fields, sa.depth), format, args...)
}

func (sa *adapter) WarningWithFields(fields LogFields, format string, args ...interface{}) {
	sa.logger.LogWithFields(LevelWarning, addCaller(fields, sa.depth), format, args...)
}

func (sa *adapter) ErrorWithFields(fields LogFields, format string, args ...interface{}) {
	sa.logger.LogWithFields(LevelError, addCaller(fields, sa.depth), format, args...)
}

func (sa *adapter) FatalWithFields(fields LogFields, format string, args ...interface{}) {
	sa.logger.LogWithFields(LevelFatal, addCaller(fields, sa.depth), format, args...)
}
