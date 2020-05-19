// Package logrus provides a Logur adapter for Logrus.
package logrus

import (
	"context"

	"github.com/sirupsen/logrus"
	"logur.dev/logur"
)

// Logger is a Logur adapter for Logrus.
type Logger struct {
	entry *logrus.Entry
}

// New returns a new Logur logger.
// If logger is nil, a default instance is created.
func New(logger *logrus.Logger) *Logger {
	if logger == nil {
		return NewFromEntry(nil)
	}

	return NewFromEntry(logrus.NewEntry(logger))
}

// NewFromEntry returns a new Logur logger from a Logrus entry.
// If entry is nil, a default instance is created.
func NewFromEntry(entry *logrus.Entry) *Logger {
	if entry == nil {
		entry = logrus.NewEntry(logrus.StandardLogger())
	}

	return &Logger{
		entry: entry,
	}
}

// Trace implements the Logur Logger interface.
func (l *Logger) Trace(msg string, fields ...map[string]interface{}) {
	if !l.entry.Logger.IsLevelEnabled(logrus.TraceLevel) {
		return
	}

	var entry = l.entry
	if len(fields) > 0 {
		entry = entry.WithFields(fields[0])
	}

	entry.Trace(msg)
}

// Debug implements the Logur Logger interface.
func (l *Logger) Debug(msg string, fields ...map[string]interface{}) {
	if !l.entry.Logger.IsLevelEnabled(logrus.DebugLevel) {
		return
	}

	var entry = l.entry
	if len(fields) > 0 {
		entry = entry.WithFields(fields[0])
	}

	entry.Debug(msg)
}

// Info implements the Logur Logger interface.
func (l *Logger) Info(msg string, fields ...map[string]interface{}) {
	if !l.entry.Logger.IsLevelEnabled(logrus.InfoLevel) {
		return
	}

	var entry = l.entry
	if len(fields) > 0 {
		entry = entry.WithFields(fields[0])
	}

	entry.Info(msg)
}

// Warn implements the Logur Logger interface.
func (l *Logger) Warn(msg string, fields ...map[string]interface{}) {
	if !l.entry.Logger.IsLevelEnabled(logrus.WarnLevel) {
		return
	}

	var entry = l.entry
	if len(fields) > 0 {
		entry = entry.WithFields(fields[0])
	}

	entry.Warn(msg)
}

// Error implements the Logur Logger interface.
func (l *Logger) Error(msg string, fields ...map[string]interface{}) {
	if !l.entry.Logger.IsLevelEnabled(logrus.ErrorLevel) {
		return
	}

	var entry = l.entry
	if len(fields) > 0 {
		entry = entry.WithFields(fields[0])
	}

	entry.Error(msg)
}

func (l *Logger) TraceContext(_ context.Context, msg string, fields ...map[string]interface{}) {
	l.Trace(msg, fields...)
}

func (l *Logger) DebugContext(_ context.Context, msg string, fields ...map[string]interface{}) {
	l.Debug(msg, fields...)
}

func (l *Logger) InfoContext(_ context.Context, msg string, fields ...map[string]interface{}) {
	l.Info(msg, fields...)
}

func (l *Logger) WarnContext(_ context.Context, msg string, fields ...map[string]interface{}) {
	l.Warn(msg, fields...)
}

func (l *Logger) ErrorContext(_ context.Context, msg string, fields ...map[string]interface{}) {
	l.Error(msg, fields...)
}

// LevelEnabled implements the Logur LevelEnabler interface.
func (l *Logger) LevelEnabled(level logur.Level) bool {
	switch level {
	case logur.Trace:
		return l.entry.Logger.IsLevelEnabled(logrus.TraceLevel)
	case logur.Debug:
		return l.entry.Logger.IsLevelEnabled(logrus.DebugLevel)
	case logur.Info:
		return l.entry.Logger.IsLevelEnabled(logrus.InfoLevel)
	case logur.Warn:
		return l.entry.Logger.IsLevelEnabled(logrus.WarnLevel)
	case logur.Error:
		return l.entry.Logger.IsLevelEnabled(logrus.ErrorLevel)
	}

	return true
}
