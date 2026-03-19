package errors

import (
	"fmt"
	"runtime"
)

// Frame represents a single stack frame.
type Frame struct {
	Function string
	File     string
	Line     int
}

func (f Frame) String() string {
	return fmt.Sprintf("%s\n\t%s:%d", f.Function, f.File, f.Line)
}

// callers captures up to 32 stack frames, skipping the given number of frames.
func callers(skip int) []Frame {
	pcs := make([]uintptr, 32)
	n := runtime.Callers(skip+2, pcs)
	if n == 0 {
		return nil
	}

	frames := runtime.CallersFrames(pcs[:n])
	result := make([]Frame, 0, n)
	for {
		f, more := frames.Next()
		result = append(result, Frame{
			Function: f.Function,
			File:     f.File,
			Line:     f.Line,
		})
		if !more {
			break
		}
	}
	return result
}
