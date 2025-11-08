package errors

import (
	"fmt"
	"runtime"
)

type Frame struct {
	Function string `json:"function"`
	Location string `json:"location"`
}

func generateFrame(pc uintptr) Frame {
	fn := runtime.FuncForPC(pc - 1)
	file, line := fn.FileLine(pc - 1)
	return Frame{
		Function: fn.Name(),
		Location: fmt.Sprintf("%s:%d", file, line),
	}
}

func generateFrames(stack []uintptr) []Frame {
	frames := make([]Frame, 0, len(stack))
	for _, pc := range stack {
		frames = append(frames, generateFrame(pc))
	}
	return frames
}
