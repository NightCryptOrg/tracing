package test

import (
	"log/slog"
	"os"
	"testing"

	"github.com/NightCryptOrg/tracing"
)

var (
	colorLogger = slog.New(tracing.NewHandler(os.Stdout, nil))
	logger      = slog.New(tracing.NewHandler(os.Stdout, &tracing.Options{
		MinLevel: slog.LevelDebug}))
)

func TestMain(m *testing.M) {

	var exitCode int
	defer func() {
		os.Exit(exitCode)
	}()

	exitCode = m.Run()
}
