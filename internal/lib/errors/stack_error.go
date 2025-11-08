package errors

import (
	"encoding/json"
	"fmt"
	"strings"
)

type stackError struct {
	err   error
	stack []uintptr
}

func (e *stackError) Error() string {
	return e.err.Error()
}

func (e *stackError) Unwrap() error {
	return e.err
}

func (e *stackError) Format(s fmt.State, verb rune) {
	if verb == 'v' && s.Flag('+') {
		var stackBuf strings.Builder
		for _, pc := range e.stack {
			f := generateFrame(pc)
			stackBuf.WriteString(fmt.Sprintf("%s\t%s\n", f.Function, f.Location))
		}
		_, _ = s.Write([]byte(fmt.Sprintf("%s\n%s", e.Error(), stackBuf.String())))
		return
	}
	_, _ = s.Write([]byte(e.Error()))
}

func (e *stackError) MarshalJSON() ([]byte, error) {
	v := struct {
		Message string  `json:"message"`
		Frames  []Frame `json:"frames"`
	}{
		Message: e.Error(),
		Frames:  generateFrames(e.stack),
	}
	return json.Marshal(v)
}
