package log

type LogLevel int

const (
	LevelFatal LogLevel = iota
	LevelError
	LevelWarning
	LevelInfo
	LevelDebug
	LevelNone
)

var names = map[LogLevel]string{
	LevelDebug:   "debug",
	LevelInfo:    "info",
	LevelWarning: "warning",
	LevelError:   "error",
	LevelFatal:   "fatal",
}

func (l LogLevel) String() string {
	if name, ok := names[l]; ok {
		return name
	}

	return "unknown"
}

func parseLogLevel(name string) LogLevel {
	for level, candidate := range names {
		if candidate == name {
			return level
		}
	}

	return LevelNone
}
