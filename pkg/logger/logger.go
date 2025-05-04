package logger

import (
	"log"
	"os"
)

type Logger struct {
	*log.Logger
	logLevel string
}

func New(level string) *Logger {
	l := &Logger{
		Logger:   log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile),
		logLevel: level,
	}
	return l
}

func (l *Logger) Debug(msg string, args ...interface{}) {
	if l.logLevel == "debug" {
		l.Printf("[DEBUG] "+msg, args...)
	}
}

func (l *Logger) Info(msg string, args ...interface{}) {
	if l.logLevel == "debug" || l.logLevel == "info" {
		l.Printf("[INFO] "+msg, args...)
	}
}

func (l *Logger) Warn(msg string, args ...interface{}) {
	if l.logLevel == "debug" || l.logLevel == "info" || l.logLevel == "warn" {
		l.Printf("[WARN] "+msg, args...)
	}
}

func (l *Logger) Error(msg string, args ...interface{}) {
	l.Printf("[ERROR] "+msg, args...)
}

func (l *Logger) Fatal(msg string, args ...interface{}) {
	l.Fatalf("[FATAL] "+msg, args...)
}
