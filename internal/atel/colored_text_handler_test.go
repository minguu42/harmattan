package atel_test

import (
	"bytes"
	"log/slog"
	"strings"
	"testing"
	"testing/slogtest"
	"time"

	"github.com/minguu42/harmattan/internal/atel"
)

func TestColoredTextHandler(t *testing.T) {
	var buf bytes.Buffer

	newHandler := func(t *testing.T) slog.Handler {
		buf.Reset()
		return atel.NewColoredTextHandler(&buf, slog.LevelDebug, true, nil)
	}
	result := func(t *testing.T) map[string]any {
		line := strings.TrimSuffix(buf.String(), "\n")
		if line == "" {
			return map[string]any{}
		}
		return parseColoredTextLine(line)
	}
	slogtest.Run(t, newHandler, result)
}

const timeLength = len(time.TimeOnly)

func parseColoredTextLine(line string) map[string]any {
	records := map[string]any{}

	if len(line) >= timeLength {
		if _, err := time.Parse(time.TimeOnly, line[:timeLength]); err == nil {
			records[slog.TimeKey] = line[:timeLength]
			line = strings.TrimPrefix(line[timeLength:], " ")
		}
	}

	space := strings.IndexByte(line, ' ')
	if space == -1 {
		if line != "" {
			records[slog.LevelKey] = line
		}
		return records
	}
	records[slog.LevelKey] = line[:space]
	line = line[space+1:]

	tokens := strings.Split(line, " ")
	firstAttrIndex := len(tokens)
	for i, token := range tokens {
		if strings.Contains(token, "=") {
			firstAttrIndex = i
			break
		}
	}

	if message := strings.Join(tokens[:firstAttrIndex], " "); message != "" {
		records[slog.MessageKey] = message
	}
	for _, token := range tokens[firstAttrIndex:] {
		key, value, _ := strings.Cut(token, "=")
		setNestedValue(records, strings.Split(key, "."), value)
	}
	return records
}

func setNestedValue(records map[string]any, parts []string, value string) {
	cur := records
	for i := 0; i < len(parts)-1; i++ {
		part := parts[i]
		sub, ok := cur[part].(map[string]any)
		if !ok {
			sub = map[string]any{}
			cur[part] = sub
		}
		cur = sub
	}
	cur[parts[len(parts)-1]] = value
}
