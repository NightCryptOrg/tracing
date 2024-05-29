package test

import (
	"bytes"
	"context"
	"log/slog"
	"reflect"
	"testing"

	"github.com/NightCryptOrg/tracing"
)

func TestFeatures(t *testing.T) {
	buf := bytes.NewBuffer(make([]byte, 0, 1024))
	buf2 := bytes.NewBuffer(make([]byte, 0, buf.Cap()))

	logRecords := func(span context.Context, records []slog.Record, handlers [2]slog.Handler) {
		// Log reference output
		for _, record := range records {
			handlers[0].Handle(span, record)
		}

		// Log compared output
		for _, record := range records {
			handlers[1].Handle(span, record)
		}
	}

	t.Run("WithGroup", func(t *testing.T) {
		buf.Reset()
		buf2.Reset()
		records := levelRecords()
		span := tracing.SpanContext(context.Background(), t.Name())

		// Log with and without group
		ref := tracing.NewHandler(buf, nil)
		group := tracing.NewHandler(buf2, nil).WithGroup("MyGroup")
		logRecords(span, records, [2]slog.Handler{ref, group})

		// Compare output
		s1 := buf.String()
		s2 := buf2.String()
		if !reflect.DeepEqual(s1, s2) {
			t.Error("log group not ignored")
		}
	})
	t.Run("WithAttrs", func(t *testing.T) {
		buf.Reset()
		buf2.Reset()
		records := levelRecords()
		span := tracing.SpanContext(context.Background(), t.Name())

		// Log with and without attribute
		ref := tracing.NewHandler(buf, nil)
		attr := tracing.NewHandler(buf2, nil).WithAttrs([]slog.Attr{{Key: "MyNum", Value: slog.IntValue(42)}})
		logRecords(span, records, [2]slog.Handler{ref, attr})

		// Compare output
		s1 := buf.String()
		s2 := buf2.String()
		if !reflect.DeepEqual(s1, s2) {
			t.Error("log attribute not ignored")
		}
	})
}
