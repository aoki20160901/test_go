package logger

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"
)

// A very small structured logger using only the stdlib. Prints JSON lines.
type Logger struct {
	out io.Writer
	mu  sync.Mutex
}

var std *Logger

// Init initializes the package logger to write to stdout. Call once from main if needed.
func Init() {
	std = New(os.Stdout)
}

// InitWithWriter initializes the package logger to write to the provided writer.
// Useful for testing or redirecting logs to a file.
func InitWithWriter(w io.Writer) {
	std = New(w)
}

// New creates a new Logger that writes to w.
func New(w io.Writer) *Logger {
	return &Logger{out: w}
}

func ensure() {
	if std == nil {
		Init()
	}
}

func (l *Logger) log(level, msg string, kvs ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	m := map[string]interface{}{"level": level, "msg": msg}
	for i := 0; i+1 < len(kvs); i += 2 {
		k := fmt.Sprint(kvs[i])
		m[k] = kvs[i+1]
	}
	b, err := json.Marshal(m)
	if err != nil {
		// fallback to simple formatting
		fmt.Fprintf(l.out, "{\"level\":%q,\"msg\":%q,\"err\":%q}\n", level, msg, err.Error())
		return
	}
	fmt.Fprintln(l.out, string(b))
}

// Info logs an info-level message. kvs are key, value pairs.
func Info(msg string, kvs ...interface{}) {
	ensure()
	std.log("INFO", msg, kvs...)
}

// Debug logs a debug-level message.
func Debug(msg string, kvs ...interface{}) {
	ensure()
	std.log("DEBUG", msg, kvs...)
}

// Error logs an error-level message.
func Error(msg string, kvs ...interface{}) {
	ensure()
	std.log("ERROR", msg, kvs...)
}
