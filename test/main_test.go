package test

import (
	"fmt"
	"log/slog"
	"os"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/NightCryptOrg/tracing"
)

var (
	logger        = slog.New(tracing.NewHandler(os.Stdout, nil))
	noColorLogger = slog.New(tracing.NewHandler(os.Stdout, &tracing.Options{
		MinLevel: slog.LevelDebug}))

	levels = []slog.Level{slog.LevelInfo, slog.LevelDebug, slog.LevelWarn, slog.LevelError}
)

func TestMain(m *testing.M) {
	var exitCode int
	defer func() {
		os.Exit(exitCode)
	}()

	slog.SetDefault(logger)

	exitCode = m.Run()
}

// Helpers

// Convert a single word to Titlecase
func titleCase(s string) string {
	if s == "" {
		return ""
	}
	var b strings.Builder
	b.WriteString(strings.ToUpper(string(s[0])))
	if len(s) > 1 {
		b.WriteString(strings.ToLower(s[1:]))
	}

	return b.String()
}

// Get a list of slog.Records for each log level @ time.Now()
func levelRecords() []slog.Record {
	records := make([]slog.Record, len(levels))
	now := time.Now()
	for i := range records {
		level := levels[i]
		// Get program counter
		pc := make([]uintptr, 1)
		runtime.Callers(0, pc)
		// Create record
		records[i] = slog.NewRecord(now, level, fmt.Sprintf("%s test record", titleCase(level.String())), pc[0])
	}

	return records
}
