package logger

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
)

// Level controls verbosity.
type Level string

const (
	LevelTrace Level = "trace"
	LevelDebug Level = "debug"
	LevelInfo  Level = "info"
	LevelWarn  Level = "warn"
	LevelError Level = "error"
	LevelFatal Level = "fatal"
	LevelPanic Level = "panic"
)

func levelRank(level Level) int {
	switch level {
	case LevelTrace:
		return 0
	case LevelDebug:
		return 1
	case LevelInfo:
		return 2
	case LevelWarn:
		return 3
	case LevelError:
		return 4
	case LevelFatal:
		return 5
	case LevelPanic:
		return 6
	default:
		panic(fmt.Sprintf("logger: invalid level %q", level))
	}
}

type Fields map[string]any

type BasicLogger struct {
	logger *log.Logger
	prefix string
	level  Level
	fields Fields
}

// NewBasicLogger creates a BasicLogger with the provided level.
// Use logger.LevelInfo (or other constants) to control verbosity.
func NewBasicLogger(level Level) CustomLogger {
	levelRank(level)
	return &BasicLogger{
		logger: log.New(os.Stdout, "", log.LstdFlags),
		level:  level,
	}
}

func (b BasicLogger) Trace(msg string) {
	if !b.shouldLog(LevelTrace) {
		return
	}
	b.logger.Printf("[TRACE] %s %s", b.prefix, msg)
}
func (b BasicLogger) Tracef(format string, args ...any) {
	if !b.shouldLog(LevelTrace) {
		return
	}
	b.output("TRACE", format, args...)
}

func (b BasicLogger) Debug(msg string) {
	if !b.shouldLog(LevelDebug) {
		return
	}
	b.logger.Printf("[DEBUG] %s %s", b.prefix, msg)
}
func (b BasicLogger) Debugf(format string, args ...any) {
	if !b.shouldLog(LevelDebug) {
		return
	}
	b.output("DEBUG", format, args...)
}

func (b BasicLogger) Info(msg string) {
	if !b.shouldLog(LevelInfo) {
		return
	}
	b.logger.Printf("[INFO] %s %s", b.prefix, msg)
}
func (b BasicLogger) Infof(format string, args ...any) {
	if !b.shouldLog(LevelInfo) {
		return
	}
	b.output("INFO", format, args...)
}

func (b BasicLogger) Warn(msg string) {
	if !b.shouldLog(LevelWarn) {
		return
	}
	b.logger.Printf("[WARN] %s %s", b.prefix, msg)
}
func (b BasicLogger) Warnf(format string, args ...any) {
	if !b.shouldLog(LevelWarn) {
		return
	}
	b.output("WARN", format, args...)
}

func (b BasicLogger) Error(msg string) {
	if !b.shouldLog(LevelError) {
		return
	}
	b.logger.Printf("[ERROR] %s %s", b.prefix, msg)
}
func (b BasicLogger) Errorf(format string, args ...any) {
	if !b.shouldLog(LevelError) {
		return
	}
	b.output("ERROR", format, args...)
}

func (b BasicLogger) Fatal(msg string) {
	if !b.shouldLog(LevelFatal) {
		return
	}
	b.logger.Fatalf("[FATAL] %s %s", b.prefix, msg)
}
func (b BasicLogger) Fatalf(format string, args ...any) {
	if !b.shouldLog(LevelFatal) {
		return
	}
	b.output("FATAL", format, args...)
}

func (b BasicLogger) Panic(msg string) {
	if !b.shouldLog(LevelPanic) {
		return
	}
	b.logger.Panicf("[PANIC] %s %s", b.prefix, msg)
}
func (b BasicLogger) Panicf(format string, args ...any) {
	if !b.shouldLog(LevelPanic) {
		return
	}
	b.output("PANIC", format, args...)
}

func (b BasicLogger) shouldLog(level Level) bool {
	return levelRank(level) >= levelRank(b.level)
}

func (b *BasicLogger) WithFields(fields Fields) CustomLogger {
	// Create a new map to hold the combined context
	newFields := make(Fields)

	// 1. Copy existing fields from the current logger instance
	for k, v := range b.fields {
		newFields[k] = v
	}

	// 2. Add the new fields provided in the arguments
	for k, v := range fields {
		newFields[k] = v
	}

	// Return a new logger instance with the merged fields
	return &BasicLogger{
		logger: b.logger,
		prefix: b.prefix,
		level:  b.level,
		fields: newFields,
	}
}

func (b *BasicLogger) SetPrefix(prefix string) CustomLogger {
	// return a new BasicLogger that shares the same underlying log.Logger
	return &BasicLogger{
		logger: b.logger,
		prefix: prefix,
		level:  b.level,
	}
}

func (b *BasicLogger) SetLevel(level Level) CustomLogger {
	levelRank(level)
	// return a new BasicLogger with same logger and prefix but new level
	return &BasicLogger{
		logger: b.logger,
		prefix: b.prefix,
		level:  level,
	}
}

// FormatPretty formats pretty JSON using the logger instance.
func (b BasicLogger) FormatPretty(data interface{}) string {
	// determine a human-friendly type name
	var typeName string
	if data == nil {
		typeName = "<nil>"
	} else {
		t := reflect.TypeOf(data)
		if t.Kind() == reflect.Ptr {
			// prefer the element name if available
			elem := t.Elem()
			if elem.Name() != "" {
				typeName = "*" + elem.Name()
			} else {
				typeName = t.String()
			}
		} else if t.Name() != "" {
			typeName = t.Name()
		} else {
			typeName = t.String()
		}
	}

	bs, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		// fallback to default formatting when JSON marshal fails
		b.logger.Printf("%s: failed to marshal data to JSON: %v; fallback: %+v\n", typeName, err, data)
		return ""
	}

	return fmt.Sprintf("%s:\n%s\n", typeName, string(bs))
}

func (b BasicLogger) Dump(data string) {
	b.logger.Printf("[DUMP] %s", data)
}

func (b BasicLogger) formatFields() string {
	if len(b.fields) == 0 {
		return ""
	}

	var parts []string
	for k, v := range b.fields {
		parts = append(parts, fmt.Sprintf("%s=%v", k, v))
	}

	// Returns: [user_id=123 ip=1.1.1.1]
	return "[" + strings.Join(parts, " ") + "] "
}

func (b BasicLogger) output(logType, format string, args ...any) {
	fieldStr := b.formatFields()

	// The format includes: [LEVEL] Prefix [Fields] Message
	newFormat := "[" + logType + "] %s %s" + format
	newArgs := append([]any{b.prefix, fieldStr}, args...)

	b.logger.Printf(newFormat, newArgs...)
}
