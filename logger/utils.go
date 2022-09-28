package logger

import (
	"path/filepath"
	"runtime"
	"strings"
)

// FunctionInfo returns the name of the function and file, the line number on the calling goroutine's stack.
// The argument skip is the number of stack frames to ascend.
func FunctionInfo(skip int) (string, string, int) {
	pc, file, line, _ := runtime.Caller(skip)
	nameEnd := filepath.Ext(runtime.FuncForPC(pc).Name())
	return strings.TrimPrefix(nameEnd, "."), file, line
}
