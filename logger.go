package log

type (
	Logger interface {
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

const (
	LevelFatal LogLevel = iota
	LevelError
	LevelWarning
	LevelInfo
	LevelDebug
)

const (
	ConsoleTimeFormat = "2006-01-02 15:04:05.000"
	JSONTimeFormat    = "2006-01-02T15:04:05.000-0700"
)
