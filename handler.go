package tracing

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"sync"
)

// Slog.handler
//
// ref https://github.com/lmittmann/tint/blob/c4b42929a2f81f8fc100808a5e28956f33fe2739/handler.go
type handler struct {
	w     io.Writer
	mutex sync.Mutex // output mutex, acquire lock before writing

	opts Options
}

func (h *handler) Enabled(_ context.Context, level slog.Level) bool {
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

func (h *handler) writeColor(buf *bytes.Buffer, color string) {
	if !h.opts.Color {
		return
	}
	buf.WriteString(color)
}

func (h *handler) writeAttr(buf *bytes.Buffer, key string, value string, first *bool) {
	if *first {
		buf.WriteRune(' ')
		*first = false
	} else {
		buf.WriteString(", ")
	}

	// Attr name
	h.writeColor(buf, ansiBold)
	buf.WriteString(key)
	h.writeColor(buf, ansiReset)
	buf.WriteString(": ")

	// Attr value
	buf.WriteString(value)
}

func (h *handler) Handle(ctx context.Context, r slog.Record) error {
	const indent = "  "

	// Grab a line buffer from the pool
	buf := GetBuffer()
	defer CloseBuffer(buf)

	// Time
	buf.WriteString(r.Time.Format(h.opts.TimeFormat))

	// Log level
	logLevel := r.Level.String()
	for i := len(logLevel); i < 6; i++ { // Right-pad loglevel to keep alignment
		buf.WriteRune(' ')
	}
	h.writeColor(buf, getLogColor(r.Level))
	buf.WriteString(r.Level.String())
	// Message
	buf.WriteString(indent)
	buf.WriteString(r.Message)
	h.writeColor(buf, ansiReset)
	buf.WriteRune('\n')

	// Message attributes
	if r.NumAttrs() > 0 {
		buf.WriteString(indent)
		h.writeColor(buf, ansiItal)
		buf.WriteString("with")
		h.writeColor(buf, ansiReset)
		first := true
		r.Attrs(func(attr slog.Attr) bool {
			h.writeAttr(buf, attr.Key, fmt.Sprint(attr.Value), &first)
			return true
		})
		buf.WriteRune('\n')
	}

	// Span(s) from innermost to outermost
	if span, ok := ctx.Value(spanKey).(*Span); ok && span != nil {
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
			first := true // skip comma separator for first field
			for name, val := range span.Fields {
				h.writeAttr(buf, name, val, &first)
			}

			buf.WriteRune('\n')
		}
	}

	h.mutex.Lock()
	buf.WriteTo(h.w)
	h.mutex.Unlock()

	return nil
}

// WithAttrs
// NOTE: This handler ignores attributes. The correct way to log
// key/value pairs is to use spans or message-level attributes. See SpanContext() for more info.
func (h *handler) WithAttrs(_ []slog.Attr) slog.Handler {
	return h
}

// WithGroup
// NOTE: This handler currently ignores groups.
func (h *handler) WithGroup(name string) slog.Handler {
	return h
}

// NewHandler - Return a slog.Handler to be used for tracing
func NewHandler(w io.Writer, opts *Options) slog.Handler {
	// Allow calling with (w, nil) to use default options
	if opts == nil {
		opts = defaultOptions()
	} else {
		// Use defaults for unspecified optional fields
		defaults := defaultOptions()
		if opts.TimeFormat == "" {
			opts.TimeFormat = defaults.TimeFormat
		}
	}
	return &handler{
		w:     w,
		mutex: sync.Mutex{},
		opts:  *opts}
}
