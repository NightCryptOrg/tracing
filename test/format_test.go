package test

import (
	"context"
	"fmt"
	"log/slog"
	"math"
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
		span := tracing.NewSpan(t.Name())
		testFormat(t, span, logger)
	})
	t.Run("No Color", func(t *testing.T) {
		span := tracing.NewSpan(t.Name())
		testFormat(t, span, noColorLogger)
	})

	t.Run("Multiple Spans", func(t *testing.T) {
		const depth = 5

		span := tracing.NewSpan(t.Name())
		for i := 0; i < depth; i++ {
			span = tracing.NewSpanCtx(span, fmt.Sprintf("Inner Span %d", i))
		}
		testFormat(t, span, logger)
		t.Run("No Color", func(t *testing.T) {
			span := tracing.NewSpan(t.Name())
			for i := 0; i < depth; i++ {
				span = tracing.NewSpanCtx(span, fmt.Sprintf("Inner Span %d", i))
			}
			testFormat(t, span, noColorLogger)
		})
	})

	t.Run("Message Attrs", func(t *testing.T) {
		args := []any{"str", "test", "int", 999, "float", math.Pi}
		const testMsg = "Test Message Attrs"
		logger.InfoContext(context.Background(), "Test Message Attrs", args...)

		t.Run("With Span", func(t *testing.T) {
			span := tracing.NewSpanCtx(context.Background(), "Message Attrs")
			logger.InfoContext(span, "Test Message Attrs w/ Span", args...)
		})
	})
	t.Run("Shorthand Methods", func(t *testing.T) {
		tracing.Info("Info message")
		tracing.Debug("Debug message")
		tracing.Warn("Warn message")
		tracing.Error("Error message")
	})
}
