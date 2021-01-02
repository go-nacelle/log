package log

type NilLogger struct{}

func NewNilLogger() Logger {
	return FromMinimalLogger(&NilLogger{})
}

func (n *NilLogger) WithFields(LogFields) MinimalLogger {
	return n
}

func (n *NilLogger) LogWithFields(LogLevel, LogFields, string, ...interface{}) {
}

func (n *NilLogger) Sync() error {
	return nil
}
