package tracing

import (
	"bytes"
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

	opts Options
}

func (h *slogHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.opts.MinLevel
}

// Log formatting colors - based on Rust's phenomenal `tracing` crate
const (
	ansiReset = "\x1b[0m"
	ansiBold  = "\x1b[1m"
	ansiItal  = "\x1b[3m"
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

func (h *slogHandler) Handle(ctx context.Context, r slog.Record) error {
	const indent = "  "

	// Grab a line buffer from the pool
	buf := GetBuffer()
	defer CloseBuffer(buf)

	// Time
	buf.WriteString(r.Time.Format(h.opts.TimeFormat))

	// Log level
	buf.WriteRune(' ')
	h.writeColor(buf, getLogColor(r.Level))
	buf.WriteString(r.Level.String())
	// Message
	buf.WriteString(indent)
	buf.WriteString(r.Message)
	h.writeColor(buf, ansiReset)
	buf.WriteRune('\n')

	// Span(s) from innermost to outermost
	if span, ok := ctx.Value(spanKey).(*Span); ok {
		for ; span != nil; span = span.parent {
			buf.WriteString(indent)
			// 'in'/'with' keywords are italicized
			// span names and field names are bolded
			h.writeColor(buf, ansiItal)
			buf.WriteString("in")
			h.writeColor(buf, ansiReset)
			buf.WriteRune(' ')

			h.writeColor(buf, ansiBold)
			buf.WriteString(span.Name)
			h.writeColor(buf, ansiReset)
			if len(span.Fields) == 0 {
				buf.WriteRune('\n')
				continue
			}

			// Span fields
			buf.WriteRune(' ')
			h.writeColor(buf, ansiItal)
			buf.WriteString("with")
			h.writeColor(buf, ansiReset)
			first := true // skip comma separator for first filed
			for name, val := range span.Fields {
				if first {
					buf.WriteRune(' ')
					first = false
				} else {
					buf.WriteString(", ")
				}

				// Field name
				h.writeColor(buf, ansiBold)
				buf.WriteString(name)
				h.writeColor(buf, ansiReset)
				buf.WriteString(": ")

				// Field value
				buf.WriteString(val)
			}

			buf.WriteRune('\n')
		}
	}

	h.mutex.Lock()
	buf.WriteTo(h.w)
	h.mutex.Unlock()

	return nil
}

func (h *slogHandler) writeColor(buf *bytes.Buffer, color string) {
	if !h.opts.Color {
		return
	}
	buf.WriteString(color)
}
