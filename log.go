package tracing

import "log/slog"

func Debug(msg string, args ...any) {
	slog.DebugContext(NewSpan(""), msg, args...)
}
func Info(msg string, args ...any) {
	slog.InfoContext(NewSpan(""), msg, args...)
}
func Warn(msg string, args ...any) {
	slog.WarnContext(NewSpan(""), msg, args...)
}
func Error(msg string, args ...any) {
	slog.ErrorContext(NewSpan(""), msg, args...)
}
