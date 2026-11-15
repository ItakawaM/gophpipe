package logging

import (
	"context"
	"io"
	"log/slog"
	"os"
)

// HandlerType is a type of slog.Handler.
type HandlerType uint8

const (
	// TextHandler logs structured plaintext logs.
	TextHandler HandlerType = iota

	// JSONHandler logs structured JSON logs.
	JSONHandler
)

const (
	// LevelFatal is a custom slog.Level representing
	// a fatal error.
	LevelFatal = slog.Level(12)
)

// levelNames maps custom slog.Level to their text.
var levelNames = map[slog.Level]string{
	LevelFatal: "FATAL",
}

func replaceAttrs(_ []string, a slog.Attr) slog.Attr {
	if a.Key == slog.LevelKey {
		level, ok := a.Value.Any().(slog.Level)
		if ok {
			levelLabel, exists := levelNames[level]
			if !exists {
				levelLabel = level.String()
			}

			a.Value = slog.StringValue(levelLabel)
		}
	}

	return a
}

// Logger is a wrapper around slog.Logger
// with added methods for Fatal logging.
type Logger struct {
	*slog.Logger

	// exit is used to test Fatal and FatalContext
	exit func(code int)
}

// New returns a new Logger configured to write structured logs.
func New(out io.Writer, t HandlerType, minLevel slog.Level) *Logger {
	var handler slog.Handler

	switch t {
	case TextHandler:
		handler = slog.NewTextHandler(out, &slog.HandlerOptions{
			Level:       minLevel,
			ReplaceAttr: replaceAttrs,
		})

	case JSONHandler:
		handler = slog.NewJSONHandler(out, &slog.HandlerOptions{
			Level:       minLevel,
			ReplaceAttr: replaceAttrs,
		})

	default:
		panic("invalid handler type")
	}

	return &Logger{
		Logger: slog.New(handler),
		exit:   os.Exit,
	}
}

// Fatal logs at LevelFatal.
// It uses context.Background internally; to specify the context, use Logger.FatalContext.
func (l *Logger) Fatal(msg string, args ...any) {
	l.FatalContext(context.Background(), msg, args...)
}

// FatalContext logs at LevelFatal with the given context.
func (l *Logger) FatalContext(ctx context.Context, msg string, args ...any) {
	l.Log(ctx, LevelFatal, msg, args...)
	l.exit(1)
}
