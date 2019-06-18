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

	LogLevel int
)

const (
	LevelFatal LogLevel = iota
	LevelError
	LevelWarning
	LevelInfo
	LevelDebug

	ConsoleTimeFormat = "2006-01-02 15:04:05.000"
	JSONTimeFormat    = "2006-01-02T15:04:05.000-0700"
)

func (l LogLevel) String() string {
	switch l {
	case LevelDebug:
		return "debug"
	case LevelInfo:
		return "info"
	case LevelWarning:
		return "warning"
	case LevelError:
		return "error"
	case LevelFatal:
		return "fatal"
	default:
		return "unknown"
	}
}
