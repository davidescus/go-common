package logger

type CustomLogger interface {
	// levels
	Trace(msg string)
	Tracef(format string, args ...any)
	Debug(msg string)
	Debugf(format string, args ...any)
	Info(msg string)
	Infof(format string, args ...any)
	Warn(msg string)
	Warnf(format string, args ...any)
	Error(msg string)
	Errorf(format string, args ...any) // <--- Add this
	Fatal(msg string)
	Fatalf(format string, args ...any)
	Panic(msg string)
	Panicf(format string, args ...any)

	WithFields(fields Fields) CustomLogger

	// SetPrefix returns a new logger instance with the given prefix
	SetPrefix(prefix string) CustomLogger
	// SetLevel returns a new logger instance with the given logging level
	SetLevel(level Level) CustomLogger

	// Debug utilities
	// FormatPretty returns a human-readable representation of the provided
	// value. Implementations should return a pretty-printed (indented)
	// JSON representation prefixed with the concrete type name when
	// possible. If formatting fails, the implementation may log an error
	// and return the empty string.
	//
	// The returned string is intended for inclusion in log messages; callers
	// should handle an empty return value (formatting failure) appropriately.
	FormatPretty(data interface{}) string

	// Dump writes strings directlly to the logger.
	Dump(data string)
}
