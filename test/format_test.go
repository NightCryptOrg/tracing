package test

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"testing"

	"github.com/NightCryptOrg/tracing"
)

func TestColors(t *testing.T) {
	levels := []slog.Level{slog.LevelInfo, slog.LevelDebug, slog.LevelWarn, slog.LevelError}

	span := tracing.SpanContext(context.Background(), t.Name())
	for _, level := range levels {
		var name strings.Builder
		lvl := level.String()
		name.WriteString(string(lvl[0]))
		name.WriteString(strings.ToLower(lvl[1:]))
		colorLogger.Log(span, level, fmt.Sprintf("%s message", name.String()))
	}
}
