package log

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

var (
	logger   Logger
	levelVar slog.LevelVar

	textEnabled bool
	jsonEnabled bool
)

type Level slog.Level

const (
	LevelDebug Level = -4
	LevelInfo  Level = 0
	LevelWarn  Level = 4
	LevelError Level = 8
)

func init() {
	SetTextLogger(os.Stdout, true)
}

func setDefaultSlogHandlerOptions(l *slog.HandlerOptions) {
	l.AddSource = true
	l.Level = &levelVar
}

// EnableTextLogger enables text logger.
func EnableTextLogger() {
	textEnabled = true
}

// EnableJsonLogger enables json logger.
func EnableJsonLogger() {
	jsonEnabled = true
}

// DisableTextLogger disables text logger.
func DisableTextLogger() {
	if !jsonEnabled {
		return
	}
	textEnabled = false
}

// DisableJsonLogger disables json logger.
func DisableJsonLogger() {
	if !textEnabled {
		return
	}
	jsonEnabled = false
}

// Default returns the default logger.
func Default() *Logger {
	return &logger
}

// AddSource adds source to slog handler options.
func AddSource(options *slog.HandlerOptions) {
	options.AddSource = true
	options.ReplaceAttr = func(groups []string, a slog.Attr) slog.Attr {
		// Remove the directory from the source's filename.
		if a.Key == slog.SourceKey {
			source := a.Value.Any().(*slog.Source)
			source.File = filepath.Base(source.File)
		}
		return a
	}
}

// SetTextLogger sets and enables text logger.
func SetTextLogger(writer io.Writer, addSource bool) {
	slogHandlerOptions := &slog.HandlerOptions{}
	setDefaultSlogHandlerOptions(slogHandlerOptions)
	if addSource {
		AddSource(slogHandlerOptions)
	}

	logger.text = slog.New(slog.NewTextHandler(writer, slogHandlerOptions))
	textEnabled = true
}

// SetJsonLogger sets and enables json logger.
func SetJsonLogger(writer io.Writer, addSource bool) {
	slogHandlerOptions := &slog.HandlerOptions{}
	setDefaultSlogHandlerOptions(slogHandlerOptions)
	if addSource {
		AddSource(slogHandlerOptions)
	}

	logger.json = slog.New(slog.NewJSONHandler(writer, slogHandlerOptions))
	jsonEnabled = true
}

// SetLevelDebug sets the default logger's level to debug.
func SetLevelDebug() {
	levelVar.Set(slog.LevelDebug)
}

// SetLevelInfo sets the default logger's level to info.
func SetLevelInfo() {
	levelVar.Set(slog.LevelInfo)
}

// SetLevelWarn sets the default logger's level to warn.
func SetLevelWarn() {
	levelVar.Set(slog.LevelWarn)
}

// SetLevelError sets the default logger's level to error.
func SetLevelError() {
	levelVar.Set(slog.LevelError)
}

func Debug(msg string, args ...any) {
	r := newRecord(slog.LevelDebug, msg)
	r.Add(args...)
	handle(nil, r, slog.LevelDebug)
}

func Info(msg string, args ...any) {
	r := newRecord(slog.LevelInfo, msg)
	r.Add(args...)
	handle(nil, r, slog.LevelInfo)
}

func Warn(msg string, args ...any) {
	r := newRecord(slog.LevelWarn, msg)
	r.Add(args...)
	handle(nil, r, slog.LevelWarn)
}

func Error(msg string, args ...any) {
	r := newRecord(slog.LevelError, msg)
	r.Add(args...)
	handle(nil, r, slog.LevelError)
}

// Debugf logs and formats a debug message.
func Debugf(format string, args ...any) {
	r := newRecord(slog.LevelDebug, format, args...)
	handle(nil, r, slog.LevelDebug)
}

// Infof logs and formats an info message.
func Infof(format string, args ...any) {
	r := newRecord(slog.LevelInfo, format, args...)
	handle(nil, r, slog.LevelInfo)
}

// Warnf logs and formats a warn message.
func Warnf(format string, args ...any) {
	r := newRecord(slog.LevelWarn, format, args...)
	handle(nil, r, slog.LevelWarn)
}

// Errorf logs and formats an error message.
func Errorf(format string, args ...any) {
	r := newRecord(slog.LevelError, format, args...)
	handle(nil, r, slog.LevelError)
}

func newRecord(level slog.Level, format string, args ...any) slog.Record {
	var pcs [1]uintptr
	runtime.Callers(3, pcs[:]) // skip [runtime.Callers, this function, this function's caller]
	if args == nil {
		return slog.NewRecord(time.Now(), level, format, pcs[0])
	}
	return slog.NewRecord(time.Now(), level, fmt.Sprintf(format, args...), pcs[0])
}

func handle(l *Logger, r slog.Record, level slog.Level) {
	if l == nil {
		if textEnabled && logger.text.Enabled(nil, level) {
			_ = logger.text.Handler().Handle(context.Background(), r)
		}

		if jsonEnabled && logger.json.Enabled(nil, level) {
			_ = logger.json.Handler().Handle(context.Background(), r)
		}
	} else {
		if textEnabled && l.text.Enabled(nil, level) {
			_ = l.text.Handler().Handle(context.Background(), r)
		}

		if jsonEnabled && l.json.Enabled(nil, level) {
			_ = l.json.Handler().Handle(context.Background(), r)
		}
	}

}
