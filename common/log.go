package common

import (
	"fmt"
	"io"
	"os"
)

type Log interface {
	Debug(message string)
	Debugf(format string, args ...any)
	Info(message string)
	Infof(format string, args ...any)
	Error(message string)
	Errorf(format string, args ...any)
	SetCurrentLogLevel(level LogLevel)
}

type LogLevel uint8

const (
	LOG_LEVEL_TRACE LogLevel = 5
	LOG_LEVEL_DEBUG LogLevel = 4
	LOG_LEVEL_INFO  LogLevel = 3
	LOG_LEVEL_WARN  LogLevel = 2
	LOG_LEVEL_ERROR LogLevel = 1
	LOG_LEVEL_NONE  LogLevel = 0
)

type LogImpl struct {
	output          io.Writer
	currentLogLevel LogLevel
}

var rootLogger = LogImpl{output: os.Stdout, currentLogLevel: LOG_LEVEL_INFO}

var _ Log = (*LogImpl)(nil)

func (l *LogImpl) SetCurrentLogLevel(level LogLevel) {
	l.currentLogLevel = level
}

func (l *LogImpl) log(message string) {
	_, _ = l.output.Write([]byte(message + "\n"))
}

func (l *LogImpl) Error(message string) {
	if l.currentLogLevel >= LOG_LEVEL_ERROR {
		l.log(message)
	}
}

func (l *LogImpl) Errorf(format string, args ...any) {
	if l.currentLogLevel >= LOG_LEVEL_ERROR {
		l.log(fmt.Sprintf(format, args...))
	}
}

func (l *LogImpl) Debug(message string) {
	if l.currentLogLevel >= LOG_LEVEL_DEBUG {
		l.log(message)
	}
}

func (l *LogImpl) Debugf(format string, args ...any) {
	if l.currentLogLevel >= LOG_LEVEL_DEBUG {
		l.log(fmt.Sprintf(format, args...))
	}
}

func (l *LogImpl) Info(message string) {
	if l.currentLogLevel >= LOG_LEVEL_INFO {
		l.log(message)
	}
}

func (l *LogImpl) Infof(format string, args ...any) {
	if l.currentLogLevel >= LOG_LEVEL_INFO {
		l.log(fmt.Sprintf(format, args...))
	}
}

func RootLogger() Log {
	return &rootLogger
}
