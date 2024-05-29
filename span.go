package tracing

import (
	"context"
	"fmt"
)

type skey string // Span ctx key

const spanKey = skey("span")

// Span - A tracing span. See https://docs.rs/tracing/latest/tracing/#spans for details
type Span struct {
	Name   string
	Fields map[string]string

	parent *Span // Parent/outer span. nil if this is the outermost span
}

// Enter a new tracing span, returning a new context with the new span.
// If the parent ctx contains a span, the new span will be added as a child
func NewSpan(parent context.Context, name string, args ...any) context.Context {
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
