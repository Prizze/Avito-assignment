package logger

//go:generate mockgen -source logger_interface.go -destination=mocks/mock_logger.go -package=mocks

type LoggerFields map[string]interface{}

type Logger interface {
	Info(msg string)
	Warn(msg string)
	Error(msg string)
	Debug(msg string)
	WithFields(fields LoggerFields) Logger
}
