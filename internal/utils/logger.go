package utils

import (
	"fmt"
	"log"
	"os"
	"time"
)

type LogLevel int

const (
	INFO LogLevel = iota
	WARNING
	ERROR
)

type Logger struct {
	infoLogger    *log.Logger
	warningLogger *log.Logger
	errorLogger   *log.Logger
}

func NewLogger() *Logger {
	flags := log.Lmsgprefix

	return &Logger{
		infoLogger:    log.New(os.Stdout, "INFO    ", flags),
		warningLogger: log.New(os.Stdout, "WARNING ", flags),
		errorLogger:   log.New(os.Stderr, "ERROR   ", flags),
	}
}

func (l *Logger) log(level LogLevel, format string, v ...interface{}) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	message := fmt.Sprintf(format, v...)
	formatted := fmt.Sprintf("[%s] %s", timestamp, message)

	switch level {
	case INFO:
		l.infoLogger.Println(formatted)
	case WARNING:
		l.warningLogger.Println(formatted)
	case ERROR:
		l.errorLogger.Println(formatted)
	}
}

func (l *Logger) Info(format string, v ...interface{}) {
	l.log(INFO, format, v...)
}

func (l *Logger) Warning(format string, v ...interface{}) {
	l.log(WARNING, format, v...)
}

func (l *Logger) Error(format string, v ...interface{}) {
	l.log(ERROR, format, v...)
}
