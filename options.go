package tracing

import "log/slog"

// Options - Settings for a Handler
type Options struct {
	MinLevel   slog.Level // Minimum level to log
	TimeFormat string     // Time format string compatible with time.Time.Format
	Color      bool       // Whether to use ANSI color sequences in output
}

func defaultOptions() *Options {
	return &Options{
		MinLevel:   slog.LevelDebug,
		TimeFormat: "02 January 2006 15:04:05.000",
		Color:      true}
}
