package log

type NilShim struct{}

func NewNilLogger() Logger {
	return adaptShim(&NilShim{})
}

func (n *NilShim) WithFields(LogFields) logShim                              { return n }
func (n *NilShim) LogWithFields(LogLevel, LogFields, string, ...interface{}) {}
func (n *NilShim) Sync() error                                               { return nil }
