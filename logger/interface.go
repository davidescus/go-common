package logger

type CustomLogger interface {
	// levels
	Trace(msg string)
	Debug(msg string)
	Info(msg string)
	Warn(msg string)
	Error(msg string)
	Fatal(msg string)
	Panic(msg string)

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
