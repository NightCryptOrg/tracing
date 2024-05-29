package test

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"testing"

	"github.com/NightCryptOrg/tracing"
)

func testFormat(t *testing.T, ctx context.Context, logger *slog.Logger) {
	levels := []slog.Level{slog.LevelInfo, slog.LevelDebug, slog.LevelWarn, slog.LevelError}

	for _, level := range levels {
		var name strings.Builder
		lvl := level.String()
		name.WriteString(string(lvl[0]))
		name.WriteString(strings.ToLower(lvl[1:]))
		logger.Log(ctx, level, fmt.Sprintf("%s message", name.String()))
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
