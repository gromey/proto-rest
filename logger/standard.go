package logger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

type format uint8

const (
	FormatText format = iota
	FormatJSON

	timeFormatDefault = "2006/01/02 15:04:05.000"
)

var _ Logger = (*StdLogger)(nil)

type Config struct {
	TimeFormat    string
	AdditionalOut io.Writer
	Format        format
	Level         Level
	FuncName      bool
}

type StdLogger struct {
	mu         sync.RWMutex
	timeFormat string
	out        io.Writer
	formatter  func(Level, string) string
	level      Level
	funcName   bool
}

func New(c *Config) *StdLogger {
	if c == nil {
		c = new(Config)
	}

	l := &StdLogger{
		timeFormat: c.TimeFormat,
		level:      c.Level,
		funcName:   c.FuncName,
	}

	if len(l.timeFormat) == 0 {
		l.timeFormat = timeFormatDefault
	}

	if c.AdditionalOut != os.Stdout && c.AdditionalOut != os.Stderr && c.AdditionalOut != nil {
		l.out = c.AdditionalOut
	}

	l.formatter = l.formatterText
	if c.Format != FormatText {
		l.formatter = l.formatterJSON
	}

	return l
}

// SetTimeFormat sets the logger time format.
func (l *StdLogger) SetTimeFormat(format string) {
	if format != "" {
		l.mu.Lock()
		defer l.mu.Unlock()
		l.timeFormat = format
	}
}

// SetFormatter sets the logger formatter.
func (l *StdLogger) SetFormatter(f format) {
	l.mu.Lock()
	defer l.mu.Unlock()
	switch f {
	case FormatJSON:
		l.formatter = l.formatterJSON
	default:
		l.formatter = l.formatterText
	}
}

// SetAdditionalOut sets an additional logger output.
func (l *StdLogger) SetAdditionalOut(out io.Writer) {
	if out != nil {
		l.mu.Lock()
		defer l.mu.Unlock()
		l.out = out
	}
}

// SetLevel sets the logger level.
func (l *StdLogger) SetLevel(level Level) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.level = level
}

// SetFuncNamePrinting sets whether the logger should print the caller function name.
func (l *StdLogger) SetFuncNamePrinting(on bool) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.funcName = on
}

func (l *StdLogger) InLevel(lvl Level) bool {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.level >= lvl
}

func (l *StdLogger) Fatalf(format string, v ...any) {
	l.printf(LevelFatal, format, v...)
	os.Exit(1)
}

func (l *StdLogger) Fatal(v ...any) {
	l.print(LevelFatal, v...)
	os.Exit(1)
}

func (l *StdLogger) Errorf(format string, v ...any) {
	if l.InLevel(LevelError) {
		l.printf(LevelError, format, v...)
	}
}

func (l *StdLogger) Error(v ...any) {
	if l.InLevel(LevelError) {
		l.print(LevelError, v...)
	}
}

func (l *StdLogger) Warnf(format string, v ...any) {
	if l.InLevel(LevelWarn) {
		l.printf(LevelWarn, format, v...)
	}
}

func (l *StdLogger) Warn(v ...any) {
	if l.InLevel(LevelWarn) {
		l.print(LevelWarn, v...)
	}
}

func (l *StdLogger) Infof(format string, v ...any) {
	if l.InLevel(LevelInfo) {
		l.printf(LevelInfo, format, v...)
	}
}

func (l *StdLogger) Info(v ...any) {
	if l.InLevel(LevelInfo) {
		l.print(LevelInfo, v...)
	}
}

func (l *StdLogger) Debugf(format string, v ...any) {
	if l.InLevel(LevelDebug) {
		l.printf(LevelDebug, format, v...)
	}
}

func (l *StdLogger) Debug(v ...any) {
	if l.InLevel(LevelDebug) {
		l.print(LevelDebug, v...)
	}
}

func (l *StdLogger) Tracef(format string, v ...any) {
	if l.InLevel(LevelTrace) {
		l.printf(LevelTrace, format, v...)
	}
}

func (l *StdLogger) Trace(v ...any) {
	if l.InLevel(LevelTrace) {
		l.print(LevelTrace, v...)
	}
}

func (l *StdLogger) printf(lvl Level, format string, v ...any) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	log := l.formatter(lvl, fmt.Sprintf(format, v...))
	if l.out != nil {
		_, _ = fmt.Fprintln(l.out, log)
	}
	_, _ = fmt.Fprintln(os.Stderr, colorWrapper(lvl, log))
}

func (l *StdLogger) print(lvl Level, v ...any) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	log := l.formatter(lvl, fmt.Sprint(v...))
	if l.out != nil {
		_, _ = fmt.Fprintln(l.out, log)
	}
	_, _ = fmt.Fprintln(os.Stderr, colorWrapper(lvl, log))

}

func colorWrapper(lvl Level, log string) string {
	switch lvl {
	case LevelFatal:
		return fmt.Sprintf("\x1b[95m%s\x1b[0m", log)
	case LevelError:
		return fmt.Sprintf("\x1b[91m%s\x1b[0m", log)
	case LevelWarn:
		return fmt.Sprintf("\x1b[93m%s\x1b[0m", log)
	case LevelInfo:
		return fmt.Sprintf("\x1b[92m%s\x1b[0m", log)
	case LevelDebug:
		return fmt.Sprintf("\x1b[94m%s\x1b[0m", log)
	default:
		return fmt.Sprintf("\x1b[96m%s\x1b[0m", log)
	}
}

func (l *StdLogger) formatterText(lvl Level, msg string) string {
	buf := new(bytes.Buffer)
	buf.WriteString(time.Now().Format(l.timeFormat))
	buf.WriteByte(' ')
	buf.WriteString(lvl.String())
	buf.WriteByte(' ')

	if l.funcName || l.level == LevelTrace {
		name, file, line := FunctionInfo(5)
		if l.level == LevelTrace {
			buf.WriteString(fmt.Sprintf("%s:%d", file, line))
			buf.WriteByte(' ')
		}
		buf.WriteString("Func: ")
		buf.WriteString(fmt.Sprintf("%s() ", name))
	}

	buf.WriteString(msg)

	return buf.String()
}

type jsonLog struct {
	Time    string `json:"time,omitempty"`
	Level   string `json:"level,omitempty"`
	File    string `json:"file,omitempty"`
	Func    string `json:"func,omitempty"`
	Message string `json:"message,omitempty"`
}

func (l *StdLogger) formatterJSON(lvl Level, msg string) string {
	buf := new(jsonLog)
	buf.Time = time.Now().Format(l.timeFormat)
	buf.Level = lvl.String()

	if l.funcName || l.level == LevelTrace {
		name, file, line := FunctionInfo(5)
		if l.level == LevelTrace {
			buf.File = fmt.Sprintf("%s:%d", file, line)
		}
		buf.Func = fmt.Sprintf("%s()", name)
	}

	buf.Message = msg

	log, _ := json.Marshal(buf)

	return string(log)
}
