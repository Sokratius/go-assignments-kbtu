package logger

import "log"

type Interface interface {
	Info(msg string, args ...any)
	Error(msg string, args ...any)
}

type Logger struct{}

func New() *Logger {
	return &Logger{}
}

func (l *Logger) Info(msg string, args ...any) {
	log.Printf("INFO: "+msg, args...)
}

func (l *Logger) Error(msg string, args ...any) {
	log.Printf("ERROR: "+msg, args...)
}
