package tracing

import (
	"context"
	"io"
	"log/slog"
	"sync"
)

// Slog.Handler Reference used
// https://github.com/lmittmann/tint/blob/c4b42929a2f81f8fc100808a5e28956f33fe2739/handler.go
type slogHandler struct {
	w       io.Writer
	mutex   sync.Mutex // line mutex
	bufPool sync.Pool  // line-buffer caching pool

	opts SlogHandlerOptions
}

func (h *slogHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.opts.MinLevel
}

// Log formatting colors - based on Rust's phenomenal `tracing` crate
const (
	ansiReset = "\x1b[0m"
	ansiBold  = "\x1b[1m"
	ansiError = "\x1b[31m"
	ansiInfo  = "\x1b[32m"
	ansiWarn  = "\x1b[33m"
	ansiDebug = "\x1b[34m"
	ansiTrace = "\x1b[35m" // TODO: see if we can implement a custom slog.Level for TRACE
)

func getLogColor(level slog.Level) string {
	switch level {
	case slog.LevelDebug.Level():
		return ansiDebug
	case slog.LevelInfo.Level():
		return ansiInfo
	case slog.LevelWarn.Level():
		return ansiWarn
	case slog.LevelError.Level():
		return ansiError
	}

	return "" // unknown level
}

func (h *slogHandler) Handle(_ context.Context, r slog.Record) error {
	// Grab a line buffer from the pool
	buf := GetBuffer()
	defer CloseBuffer(buf)

	// Write time
	buf.WriteString(r.Time.Format(h.opts))

	// Write log level and message
	buf.WriteRune(' ')
	if h.opts.Color {
		buf.WriteString(getLogColor(r.Level))
	}
	buf.WriteString(r.Level.String())
	buf.WriteString("  ")
	buf.WriteString(r.Message)
	buf.WriteString(ansiReset)

	// Write span(s) from innermost to outer
	// TODO

	h.mutex.Lock()
	buf.WriteTo(h.w)
	h.mutex.Unlock()

	return nil
}
