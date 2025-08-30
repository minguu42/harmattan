package applog_test

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"testing"
	"testing/slogtest"

	"github.com/minguu42/harmattan/internal/applog"
)

func TestJSONIndentHandler(t *testing.T) {
	var buf bytes.Buffer

	newHandler := func(t *testing.T) slog.Handler {
		buf.Reset()
		return applog.NewJSONIndentHandler(&buf, nil)
	}
	result := func(t *testing.T) map[string]any {
		line := buf.Bytes()
		if len(line) == 0 {
			return map[string]any{}
		}

		var m map[string]any
		if err := json.Unmarshal(line, &m); err != nil {
			t.Fatal(err)
		}
		return m
	}
	slogtest.Run(t, newHandler, result)
}
