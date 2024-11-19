package test

import (
	"testing"

	"github.com/NightCryptOrg/tracing"
)

func TestNaming(t *testing.T) {
	t.Run("Callers", func(t *testing.T) {
		span := tracing.NewSpan("", "str", "test", "int", 999)
		testFormat(t, span, colorLogger)
	})
}
