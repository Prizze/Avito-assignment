package logger

type LoggerFields map[string]interface{}

type Logger interface {
	Info(msg string)
	Warn(msg string)
	Error(msg string)
	Debug(msg string)
	WithFields(fields LoggerFields) Logger
}