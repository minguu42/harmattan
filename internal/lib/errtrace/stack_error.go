package errtrace

import (
	"encoding/json"
	"fmt"
	"runtime"
	"strings"
)

type StackError struct {
	err   error
	stack []uintptr
}

func (e *StackError) Error() string {
	return e.err.Error()
}

func (e *StackError) Unwrap() error {
	return e.err
}

func (e *StackError) MarshalJSON() ([]byte, error) {
	v := struct {
		Message string  `json:"message"`
		Frames  []Frame `json:"frames"`
	}{
		Message: e.Error(),
		Frames:  generateFrames(e.stack),
	}
	return json.Marshal(v)
}

type Frame struct {
	Function string `json:"function"`
	Location string `json:"location"`
}

func generateFrames(stack []uintptr) []Frame {
	frames := make([]Frame, 0, len(stack))
	for _, pc := range stack {
		fn := runtime.FuncForPC(pc - 1)
		file, line := fn.FileLine(pc - 1)

		// 不要なフレームを出力しないためにフレームワーク部分のフレームは出力しない
		if name := file[strings.LastIndex(file, "/")+1:]; name == "oas_handlers_gen.go" {
			break
		}

		frames = append(frames, Frame{
			Function: fn.Name(),
			Location: fmt.Sprintf("%s:%d", file, line),
		})
	}
	return frames
}
