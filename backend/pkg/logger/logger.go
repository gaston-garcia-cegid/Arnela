package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

// Logger wraps zerolog.Logger
type Logger struct {
	*zerolog.Logger
}

// NewLogger creates a new logger instance
func NewLogger(env string) *Logger {
	var log zerolog.Logger

	// Configure based on environment
	if env == "development" {
		// Pretty console output for development
		output := zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		}
		log = zerolog.New(output).With().Timestamp().Caller().Logger()
	} else {
		// JSON output for production
		log = zerolog.New(os.Stdout).With().Timestamp().Caller().Logger()
	}

	// Set global log level
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	return &Logger{
		Logger: &log,
	}
}

// WithContext adds context fields to the logger
func (l *Logger) WithContext(fields map[string]interface{}) *Logger {
	ctx := l.Logger.With()
	for k, v := range fields {
		ctx = ctx.Interface(k, v)
	}
	newLogger := ctx.Logger()
	return &Logger{Logger: &newLogger}
}

// Info logs an info message
func (l *Logger) Info(msg string, fields ...map[string]interface{}) {
	event := l.Logger.Info()
	if len(fields) > 0 {
		for k, v := range fields[0] {
			event = event.Interface(k, v)
		}
	}
	event.Msg(msg)
}

// Error logs an error message
func (l *Logger) Error(msg string, err error, fields ...map[string]interface{}) {
	event := l.Logger.Error().Err(err)
	if len(fields) > 0 {
		for k, v := range fields[0] {
			event = event.Interface(k, v)
		}
	}
	event.Msg(msg)
}

// Warn logs a warning message
func (l *Logger) Warn(msg string, fields ...map[string]interface{}) {
	event := l.Logger.Warn()
	if len(fields) > 0 {
		for k, v := range fields[0] {
			event = event.Interface(k, v)
		}
	}
	event.Msg(msg)
}

// Debug logs a debug message
func (l *Logger) Debug(msg string, fields ...map[string]interface{}) {
	event := l.Logger.Debug()
	if len(fields) > 0 {
		for k, v := range fields[0] {
			event = event.Interface(k, v)
		}
	}
	event.Msg(msg)
}

// Fatal logs a fatal message and exits
func (l *Logger) Fatal(msg string, err error, fields ...map[string]interface{}) {
	event := l.Logger.Fatal().Err(err)
	if len(fields) > 0 {
		for k, v := range fields[0] {
			event = event.Interface(k, v)
		}
	}
	event.Msg(msg)
}
