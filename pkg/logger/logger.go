package logger

import (
	"log"
)

type Logger struct {
	logger *log.Logger
}

func NewLogger() *Logger {
	return &Logger{
		logger: log.New(log.Writer(), "[APP] ", log.LstdFlags),
	}
}

func (l *Logger) Info(msg string, args ...interface{}) {
	l.logger.Printf("[INFO] "+msg, args...)
}

func (l *Logger) Error(msg string, args ...interface{}) {
	l.logger.Printf("[ERROR] "+msg, args...)
}

func (l *Logger) Fatal(msg string, args ...interface{}) {
	l.logger.Fatalf("[FATAL] "+msg, args...)
}
