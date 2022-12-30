package logger

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"strings"
)

type Level uint8

const (
	LevelFatal Level = iota
	LevelError
	LevelWarn
	LevelInfo
	LevelDebug
	LevelTrace
)

func (lvl Level) String() string {
	switch lvl {
	case LevelFatal:
		return "FATAL"
	case LevelError:
		return "ERROR"
	case LevelWarn:
		return "WARN"
	case LevelInfo:
		return "INFO"
	case LevelDebug:
		return "DEBUG"
	default:
		return "TRACE"
	}
}

// ParseLevel takes a string level and returns the logger Level constant.
func ParseLevel(lvl string) (Level, error) {
	switch strings.ToLower(lvl) {
	case "fatal":
		return LevelFatal, nil
	case "error":
		return LevelError, nil
	case "warn", "warning":
		return LevelWarn, nil
	case "info":
		return LevelInfo, nil
	case "debug":
		return LevelDebug, nil
	case "trace":
		return LevelTrace, nil
	default:
		return LevelFatal, fmt.Errorf("not a valid logger Level: %s", lvl)
	}
}

type Logger interface {
	InLevel(lvl Level) bool
	Fatalf(format string, v ...any)
	Fatal(v ...any)
	Errorf(format string, v ...any)
	Error(v ...any)
	Warnf(format string, v ...any)
	Warn(v ...any)
	Infof(format string, v ...any)
	Info(v ...any)
	Debugf(format string, v ...any)
	Debug(v ...any)
	Tracef(format string, v ...any)
	Trace(v ...any)
}

var std Logger = New(&Config{Level: LevelTrace})

func SetLogger(logger Logger) {
	std = logger
}

// InLevel returns true if the given level is less than or equal to the current logger level.
func InLevel(level Level) bool {
	return std.InLevel(level)
}

func Fatalf(format string, v ...any) {
	std.Fatalf(format, v...)
}

func Fatal(v ...any) {
	std.Fatal(v...)
}

func Errorf(format string, v ...any) {
	std.Errorf(format, v...)
}

func Error(v ...any) {
	std.Error(v...)
}

func Warnf(format string, v ...any) {
	std.Warnf(format, v...)
}

func Warn(v ...any) {
	std.Warn(v...)
}

func Infof(format string, v ...any) {
	std.Infof(format, v...)
}

func Info(v ...any) {
	std.Info(v...)
}

func Debugf(format string, v ...any) {
	std.Debugf(format, v...)
}

func Debug(v ...any) {
	std.Debug(v...)
}

func Tracef(format string, v ...any) {
	std.Tracef(format, v...)
}

func Trace(v ...any) {
	std.Trace(v...)
}

// DumpHttpRequest dumps the HTTP request and prints out with logFunc.
func DumpHttpRequest(r *http.Request, logFunc func(v ...any)) {
	dumpFunc := httputil.DumpRequestOut
	if r.URL.Scheme == "" || r.URL.Host == "" {
		dumpFunc = httputil.DumpRequest
	}
	b, err := dumpFunc(r, true)
	if err != nil {
		if InLevel(LevelError) {
			Error("REQUEST LOG error: ", err)
		}
		return
	}
	logFunc("REQUEST: ", string(b))
}

// DumpHttpResponse dumps the HTTP response and prints out with logFunc.
func DumpHttpResponse(r *http.Response, logFunc func(v ...any)) {
	b, err := httputil.DumpResponse(r, true)
	if err != nil {
		if InLevel(LevelError) {
			Error("RESPONSE LOG error: ", err)
		}
		return
	}
	logFunc("RESPONSE: ", string(b))
}
