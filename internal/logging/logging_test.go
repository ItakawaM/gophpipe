package logging

import (
	"bytes"
	"io"
	"log/slog"
	"strings"
	"testing"

	"github.com/ItakawaM/gophpipe/internal/testutil"
)

// TestNew tests Logger creation.
func TestNew(t *testing.T) {
	tests := []struct {
		name      string
		out       io.Writer
		t         HandlerType
		minLevel  slog.Level
		wantPanic bool
	}{
		{
			name:     "info bytes text handler",
			out:      bytes.NewBuffer(nil),
			t:        TextHandler,
			minLevel: slog.LevelInfo,
		},
		{
			name:     "warn bytes json handler",
			out:      bytes.NewBuffer(nil),
			t:        JSONHandler,
			minLevel: slog.LevelWarn,
		},
		{
			name:      "unknown handler",
			out:       bytes.NewBuffer(nil),
			t:         HandlerType(213),
			minLevel:  slog.LevelInfo,
			wantPanic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				testutil.AssertPanic(t, recover(), tt.wantPanic)
			}()

			_ = New(tt.out, tt.t, tt.minLevel)
		})
	}
}

// TestReplaceAttrs tests correct LevelKey=LevelFatal replacement.
func TestReplaceAttrs(t *testing.T) {
	tests := []struct {
		name string
		attr slog.Attr
		want slog.Attr
	}{
		{
			name: "fatal",
			attr: slog.Any(slog.LevelKey, LevelFatal),
			want: slog.Any(slog.LevelKey, "FATAL"),
		},
		{
			name: "error",
			attr: slog.Any(slog.LevelKey, slog.LevelError),
			want: slog.Any(slog.LevelKey, "ERROR"),
		},
		{
			name: "unchanged",
			attr: slog.Any(slog.MessageKey, "UNCHANGED"),
			want: slog.Any(slog.MessageKey, "UNCHANGED"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := replaceAttrs(nil, tt.attr)
			testutil.AssertDeepEqual(t, got, tt.want)
		})
	}
}

// TestFatal tests fatal logging.
func TestFatal(t *testing.T) {
	tests := []struct {
		name string
		msg  string
		args []any
		want string
	}{
		{
			name: "fatal no args",
			msg:  "something happened",
			args: nil,
		},
		{
			name: "fatal with args",
			msg:  "something happened",
			args: []any{"key", "value"},
			want: "key=value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buffer := bytes.NewBuffer(nil)
			logger := New(buffer, TextHandler, slog.LevelInfo)

			// replace exit with a dummy to not crash the test
			logger.exit = func(code int) {}

			logger.Fatal(tt.msg, tt.args...)
			got := buffer.String()
			if !strings.Contains(got, "FATAL") {
				t.Errorf("Fatal(%q, %q): got %q, want FATAL", tt.msg, tt.args, got)
			}

			if tt.want != "" && !strings.Contains(got, tt.want) {
				t.Errorf("Fatal(%q, %q): got %q, want %q", tt.msg, tt.args, got, tt.want)
			}
		})
	}
}
