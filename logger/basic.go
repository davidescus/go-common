package logger

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"reflect"
)

// Level controls verbosity. Lower values are more verbose.
type Level int

const (
	LevelTrace Level = iota
	LevelDebug
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
	LevelPanic
)

type BasicLogger struct {
	logger *log.Logger
	prefix string
	level  Level
}

// NewBasicLogger creates a BasicLogger with the provided level.
// Use logger.LevelInfo (or other constants) to control verbosity.
func NewBasicLogger(level Level) CustomLogger {
	return &BasicLogger{
		logger: log.New(os.Stdout, "", log.LstdFlags),
		level:  level,
	}
}

func (b BasicLogger) Trace(msg string) {
	if LevelTrace < b.level {
		return
	}
	b.logger.Printf("[TRACE] %s %s", b.prefix, msg)
}

func (b BasicLogger) Debug(msg string) {
	if LevelDebug < b.level {
		return
	}
	b.logger.Printf("[DEBUG] %s %s", b.prefix, msg)
}

func (b BasicLogger) Info(msg string) {
	if LevelInfo < b.level {
		return
	}
	b.logger.Printf("[INFO] %s %s", b.prefix, msg)
}

func (b BasicLogger) Warn(msg string) {
	if LevelWarn < b.level {
		return
	}
	b.logger.Printf("[WARN] %s %s", b.prefix, msg)
}

func (b BasicLogger) Error(msg string) {
	if LevelError < b.level {
		return
	}
	b.logger.Printf("[ERROR] %s %s", b.prefix, msg)
}

func (b BasicLogger) Fatal(msg string) {
	if LevelFatal < b.level {
		return
	}
	b.logger.Fatalf("[FATAL] %s %s", b.prefix, msg)
}

func (b BasicLogger) Panic(msg string) {
	if LevelPanic < b.level {
		return
	}
	b.logger.Panicf("[PANIC] %s %s", b.prefix, msg)
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
