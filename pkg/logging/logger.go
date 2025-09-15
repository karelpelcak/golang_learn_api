package logging

import (
	"log"
	"os"
)

// Logger provides a simple logging interface
type Logger struct {
	infoLogger  *log.Logger
	errorLogger *log.Logger
}

// NewLogger creates a new logger instance
func NewLogger() *Logger {
	return &Logger{
		infoLogger:  log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile),
		errorLogger: log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
	}
}

// Info logs an info message
func (l *Logger) Info(message string) {
	l.infoLogger.Println(message)
}

// Error logs an error message
func (l *Logger) Error(message string) {
	l.errorLogger.Println(message)
}