package log

type (
	Logger interface {
		WithIndirectCaller(frames int) Logger
		WithFields(LogFields) Logger
		LogWithFields(LogLevel, LogFields, string, ...interface{})
		Sync() error

		// Convenience Methods
		Debug(string, ...interface{})
		Info(string, ...interface{})
		Warning(string, ...interface{})
		Error(string, ...interface{})
		Fatal(string, ...interface{})
		DebugWithFields(LogFields, string, ...interface{})
		InfoWithFields(LogFields, string, ...interface{})
		WarningWithFields(LogFields, string, ...interface{})
		ErrorWithFields(LogFields, string, ...interface{})
		FatalWithFields(LogFields, string, ...interface{})
	}
)
