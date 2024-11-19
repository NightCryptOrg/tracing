package tracing

import (
	"context"
	"fmt"
	"runtime"
)

type skey string // Span ctx key

const spanKey = skey("span")

// Span - A tracing span. See https://docs.rs/tracing/latest/tracing/#spans for details
type Span struct {
	Name   string
	Fields map[string]string

	parent *Span // Parent/outer span. nil if this is the outermost span
}

// Enter a new tracing span, returning its context to be used with `slog` methods.
// If name == "", then runtime.Callers will be used to get the names of up to 2 functions above NewSpan in the call stack
// (and will create a context 2 levels deep with this information)
func NewSpan(name string, args ...any) context.Context {
	ctx := context.Background()
	// Auto-name using 2 parent callers
	if name == "" {
		callers := getCallers(2)
		for _, fn := range callers[:len(callers)-1] {
			ctx = NewSpanCtx(ctx, fn)
		}
		return NewSpanCtx(ctx, callers[len(callers)-1], args...)
	}

	return NewSpanCtx(ctx, name, args...)
}

// Enter a new tracing span, returning a new context with the new span.
// If the parent ctx contains a span, the new span will be added as a child
func NewSpanCtx(parent context.Context, name string, args ...any) context.Context {
	var span Span
	span.Name = name
	span.Fields = make(map[string]string)
	for i := 0; i+1 < len(args); i++ {
		key, ok := args[i].(string)
		if !ok {
			// Skip invalid arguments
			continue
		}
		val := args[i+1]

		span.Fields[key] = fmt.Sprint(val)

		// Skip to the next k/v arg pair
		i++
	}

	if outer, ok := parent.Value(spanKey).(*Span); ok {
		span.parent = outer
	}
	return context.WithValue(parent, spanKey, &span)
}

// Get the names of {n} callers of the function where getCaller is called
func getCallers(n int) (callers []string) {
	pcs := make([]uintptr, n)
	n = runtime.Callers(3, pcs) // Skip runtime.Callers, getCaller, and its calling function (i.e NewSpanCtx)
	if n == 0 {
		return callers
	}

	pcs = pcs[:n] // Limit to valid pcs
	callers = make([]string, 0, n)
	frames := runtime.CallersFrames(pcs)

	hasNext := true // whether the next frame can be iterated to
	var frame runtime.Frame
	for hasNext {
		frame, hasNext = frames.Next()
		callers = append(callers, frame.Function)
	}

	return callers
}
