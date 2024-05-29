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
}
