package helheim_go

type Logger interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(err error, msg string)
	Panic(err error, msg string)
	Fatal(err error, msg string)

	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(err error, format string, args ...interface{})
	Panicf(err error, format string, args ...interface{})
	Fatalf(err error, format string, args ...interface{})
}

type noopLogger struct {
}

func NewNoopLogger() Logger {
	return &noopLogger{}
}

func (n noopLogger) Debug(args ...interface{}) {
	return
}

func (n noopLogger) Info(args ...interface{}) {
	return
}

func (n noopLogger) Warn(args ...interface{}) {
	return
}

func (n noopLogger) Error(err error, msg string) {
	return
}

func (n noopLogger) Panic(err error, msg string) {
	return
}

func (n noopLogger) Fatal(err error, msg string) {
	return
}

func (n noopLogger) Debugf(format string, args ...interface{}) {
	return
}

func (n noopLogger) Infof(format string, args ...interface{}) {
	return
}

func (n noopLogger) Warnf(format string, args ...interface{}) {
	return
}

func (n noopLogger) Errorf(err error, format string, args ...interface{}) {
	return
}

func (n noopLogger) Panicf(err error, format string, args ...interface{}) {
	return
}

func (n noopLogger) Fatalf(err error, format string, args ...interface{}) {
	return
}
