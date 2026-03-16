package errtrace

import (
	"fmt"
	"log/slog"
	"runtime"
	"strings"
)

type StackError struct {
	err   error
	stack []uintptr
	attrs []slog.Attr
}

func (e *StackError) Attrs() []slog.Attr {
	return e.attrs
}

func (e *StackError) Error() string {
	return e.err.Error()
}

func (e *StackError) Format(s fmt.State, verb rune) {
	if verb == 'v' && s.Flag('+') {
		var buf strings.Builder
		buf.WriteString(e.Error())
		if len(e.attrs) > 0 {
			buf.WriteString(" [")
			for i, a := range e.attrs {
				if i > 0 {
					buf.WriteByte(' ')
				}
				buf.WriteString(fmt.Sprintf("%s=%v", a.Key, a.Value))
			}
			buf.WriteByte(']')
		}
		buf.WriteByte('\n')
		for i, f := range generateFrames(e.stack) {
			if i > 0 {
				buf.WriteByte('\n')
			}
			buf.WriteString(fmt.Sprintf("%s\n\t%s", f.Function, f.Location))
		}
		_, _ = s.Write([]byte(buf.String()))
		return
	}
	_, _ = s.Write([]byte(e.Error()))
}

func (e *StackError) Frames() []Frame {
	return generateFrames(e.stack)
}

func (e *StackError) Unwrap() error {
	return e.err
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
