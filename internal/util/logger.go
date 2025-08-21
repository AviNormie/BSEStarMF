package util

import (
	"log"
	"os"
)

type Logger interface {
	Info(format string, v ...interface{})
	Warn(format string, v ...interface{})
	Error(format string, v ...interface{})
	Fatal(format string, v ...interface{})
}

type StandardLogger struct {
	infoLogger  *log.Logger
	warnLogger  *log.Logger
	errorLogger *log.Logger
}

func NewStandardLogger() *StandardLogger {
	return &StandardLogger{
		infoLogger:  log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile),
		warnLogger:  log.New(os.Stdout, "WARN: ", log.Ldate|log.Ltime|log.Lshortfile),
		errorLogger: log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
	}
}

func (l *StandardLogger) Info(format string, v ...interface{}) {
	l.infoLogger.Printf(format+"\n", v...)
}

func (l *StandardLogger) Warn(format string, v ...interface{}) {
	l.warnLogger.Printf(format+"\n", v...)
}

func (l *StandardLogger) Error(format string, v ...interface{}) {
	l.errorLogger.Printf(format+"\n", v...)
}

func (l *StandardLogger) Fatal(format string, v ...interface{}) {
	l.errorLogger.Fatalf(format+"\n", v...)
}