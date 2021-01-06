package logging

import (
	"fmt"

	log "github.com/sirupsen/logrus"
)

var (
	defaultFields   map[string]string
	defaultLogLevel = "info"
)

// StructuredLogger is a logger based on logrus.
type StructuredLogger struct {
	logger *log.Logger
}

// NewStructuredLogger creates a new structured logger.
func NewStructuredLogger() *StructuredLogger {
	return &StructuredLogger{
		logger: log.New(),
	}
}

// SetLevel allows the log level to be set.
func (l StructuredLogger) SetLevel(levelName string) {
	if levelName == "" {
		levelName = defaultLogLevel
	}

	level, err := log.ParseLevel(levelName)
	if err != nil {
		l.logger.Warn(fmt.Sprintf("could not parse log level '%s', logging will proceed at %s level", levelName, defaultLogLevel))
		level, _ = log.ParseLevel(defaultLogLevel)
	}

	l.logger.SetLevel(level)
}

// LogJSON determines whether or not to format the logs as JSON.
func (l StructuredLogger) SetLogJSON(value bool) {
	if value {
		l.logger.SetFormatter(&log.JSONFormatter{})
	}
}

// SetDefaultFields sets fields to be logged on every use of the logger.
func (l StructuredLogger) SetDefaultFields(defaultFields map[string]string) {
	l.logger.AddHook(&defaultFieldHook{})
}

// Error logs an error message.
func (l StructuredLogger) Error(msg string, fields ...interface{}) {
	l.logger.WithFields(createFieldMap(fields)).Error(msg)
}

// Warn logs an warning message.
func (l StructuredLogger) Warn(msg string, fields ...interface{}) {
	l.logger.WithFields(createFieldMap(fields)).Warn(msg)
}

// Info logs an info message.
func (l StructuredLogger) Info(msg string, fields ...interface{}) {
	l.logger.WithFields(createFieldMap(fields)).Info(msg)
}

// Debug logs a debug message.
func (l StructuredLogger) Debug(msg string, fields ...interface{}) {
	l.logger.WithFields(createFieldMap(fields)).Debug(msg)
}

// Trace logs a trace message.
func (l StructuredLogger) Trace(msg string, fields ...interface{}) {
	l.logger.WithFields(createFieldMap(fields)).Trace(msg)
}

func createFieldMap(fields ...interface{}) map[string]interface{} {
	m := map[string]interface{}{}

	fields = fields[0].([]interface{})

	for i := 0; i < len(fields); i += 2 {
		m[fields[i].(string)] = fields[i+1]
	}

	return m
}

type defaultFieldHook struct{}

func (h *defaultFieldHook) Levels() []log.Level {
	return log.AllLevels
}

func (h *defaultFieldHook) Fire(e *log.Entry) error {
	for k, v := range defaultFields {
		e.Data[k] = v
	}
	return nil
}
