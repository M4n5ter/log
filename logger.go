package log

import (
	"context"
	"log/slog"
)

type Logger struct {
	text *slog.Logger
	json *slog.Logger
}

func (l *Logger) Debugf(format string, args ...any) {
	r := newRecord(slog.LevelDebug, format, args...)
	handle(l, r, slog.LevelDebug)
}

func (l *Logger) Infof(format string, args ...any) {
	r := newRecord(slog.LevelInfo, format, args...)
	handle(l, r, slog.LevelInfo)
}

func (l *Logger) Warnf(format string, args ...any) {
	r := newRecord(slog.LevelWarn, format, args...)
	handle(l, r, slog.LevelWarn)
}

func (l *Logger) Errorf(format string, args ...any) {
	r := newRecord(slog.LevelError, format, args...)
	handle(l, r, slog.LevelError)
}

// With works like slog.Logger.With
func (l *Logger) With(args ...any) *Logger {
	text := l.text.With(args...)
	json := l.json.With(args...)
	return &Logger{text: text, json: json}
}

// WithGroup works like slog.Logger.WithGroup
func (l *Logger) WithGroup(name string) *Logger {
	text := l.text.WithGroup(name)
	json := l.json.WithGroup(name)
	return &Logger{text: text, json: json}
}

// Log works like slog.Logger.Log
func (l *Logger) Log(ctx context.Context, level Level, msg string, args ...any) {
	lv := slog.Level(level)
	if ctx == nil {
		ctx = context.Background()
	}

	r := newRecord(lv, msg)
	r.Add(args...)
	if textEnabled && l.text.Enabled(ctx, lv) {
		l.text.Handler().Handle(ctx, r)
	}

	if jsonEnabled && l.json.Enabled(ctx, lv) {
		l.json.Handler().Handle(ctx, r)
	}
}

// LogAttrs works like slog.Logger.LogAttrs
func (l *Logger) LogAttrs(ctx context.Context, level Level, msg string, attrs ...Attr) {
	lv := slog.Level(level)
	if ctx == nil {
		ctx = context.Background()
	}

	r := newRecord(lv, msg)
	r.AddAttrs(attrs...)
	if textEnabled && l.text.Enabled(ctx, lv) {
		l.text.Handler().Handle(ctx, r)
	}

	if jsonEnabled && l.json.Enabled(ctx, lv) {
		l.json.Handler().Handle(ctx, r)
	}
}