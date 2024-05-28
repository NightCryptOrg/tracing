package test

import (
	"log/slog"
	"os"
	"testing"

	"github.com/NightCryptOrg/tracing"
)

var (
	colorLogger *slog.Logger
	logger      *slog.Logger
)

func TestMain(m *testing.M) {
	// Initialize loggers
	colorLogger = slog.New(tracing.NewHandler(os.Stdout, nil))

	var exitCode int
	defer func() {
		os.Exit(exitCode)
	}()

	exitCode = m.Run()
}
