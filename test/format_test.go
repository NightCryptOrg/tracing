package test

import (
	"context"
	"fmt"
	"log/slog"
	"testing"

	"github.com/NightCryptOrg/tracing"
)

func testFormat(_ *testing.T, ctx context.Context, logger *slog.Logger) {
	for _, level := range levels {
		logger.Log(ctx, level, fmt.Sprintf("%s message", titleCase(level.String())))
	}
}

func TestFormat(t *testing.T) {
	t.Run("Color", func(t *testing.T) {
		span := tracing.SpanContext(context.Background(), t.Name())
		testFormat(t, span, colorLogger)
	})
	t.Run("No Color", func(t *testing.T) {
		span := tracing.SpanContext(context.Background(), t.Name())
		testFormat(t, span, logger)
	})

	t.Run("Multiple Spans", func(t *testing.T) {
		const depth = 5

		span := tracing.SpanContext(context.Background(), t.Name())
		for i := 0; i < depth; i++ {
			span = tracing.SpanContext(span, fmt.Sprintf("Inner Span %d", i))
		}
		testFormat(t, span, colorLogger)
		t.Run("No Color", func(t *testing.T) {
			testFormat(t, span, logger)
		})
	})
}
