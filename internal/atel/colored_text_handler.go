package atel

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/minguu42/harmattan/internal/lib/errtrace"
)

type colorCodes struct {
	reset  string
	gray   string
	red    string
	green  string
	yellow string
	cyan   string
}

var ansiColors = colorCodes{
	reset:  "\033[0m",
	gray:   "\033[90m",
	red:    "\033[31m",
	green:  "\033[32m",
	yellow: "\033[33m",
	cyan:   "\033[36m",
}

type ColoredTextHandler struct {
	level       slog.Leveler
	attrsPrefix string
	groupPrefix string
	colors      colorCodes
	ignoreKeys  map[string]struct{}

	w  io.Writer
	mu *sync.Mutex
}

func NewColoredTextHandler(w io.Writer, level slog.Leveler, noColor bool, ignoreKeys []string) *ColoredTextHandler {
	var colors colorCodes
	if !noColor {
		colors = ansiColors
	}
	ignoreKeysSet := make(map[string]struct{}, len(ignoreKeys))
	for _, k := range ignoreKeys {
		ignoreKeysSet[k] = struct{}{}
	}
	return &ColoredTextHandler{
		level:      level,
		colors:     colors,
		ignoreKeys: ignoreKeysSet,
		w:          w,
		mu:         &sync.Mutex{},
	}
}

func (h *ColoredTextHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.level.Level()
}

func (h *ColoredTextHandler) Handle(_ context.Context, record slog.Record) error {
	var builder strings.Builder
	if !record.Time.IsZero() {
		fmt.Fprintf(&builder, "%s%s%s ", h.colors.gray, record.Time.Format(time.TimeOnly), h.colors.reset)
	}

	switch record.Level {
	case slog.LevelDebug:
		fmt.Fprintf(&builder, "%s%s%s", h.colors.cyan, record.Level.String(), h.colors.reset)
	case slog.LevelInfo:
		fmt.Fprintf(&builder, "%s%s%s", h.colors.green, record.Level.String(), h.colors.reset)
	case slog.LevelWarn:
		fmt.Fprintf(&builder, "%s%s%s", h.colors.yellow, record.Level.String(), h.colors.reset)
	case slog.LevelError:
		fmt.Fprintf(&builder, "%s%s%s", h.colors.red, record.Level.String(), h.colors.reset)
	default:
		builder.WriteString(record.Level.String())
	}

	if record.Message != "" {
		builder.WriteByte(' ')
		builder.WriteString(record.Message)
	}

	if h.attrsPrefix != "" {
		builder.WriteString(h.attrsPrefix)
	}

	var errAttr slog.Attr
	record.Attrs(func(attr slog.Attr) bool {
		if attr.Key == "error" {
			if resolved := attr.Value.Resolve(); resolved.Kind() == slog.KindAny {
				if _, ok := resolved.Any().(error); ok {
					errAttr = attr
					return true
				}
			}
		}
		h.appendAttr(&builder, attr, h.groupPrefix)
		return true
	})
	if err, ok := errAttr.Value.Resolve().Any().(error); ok {
		fmt.Fprintf(&builder, "\n%+v", err)
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	builder.WriteByte('\n')
	if _, err := h.w.Write([]byte(builder.String())); err != nil {
		return errtrace.Wrap(err)
	}
	return nil
}

func (h *ColoredTextHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	if len(attrs) == 0 {
		return h
	}

	var builder strings.Builder
	for _, attr := range attrs {
		h.appendAttr(&builder, attr, h.groupPrefix)
	}
	return &ColoredTextHandler{
		level:       h.level,
		attrsPrefix: h.attrsPrefix + builder.String(),
		groupPrefix: h.groupPrefix,
		colors:      h.colors,
		ignoreKeys:  h.ignoreKeys,
		w:           h.w,
		mu:          h.mu,
	}
}

func (h *ColoredTextHandler) WithGroup(name string) slog.Handler {
	if name == "" {
		return h
	}
	return &ColoredTextHandler{
		level:       h.level,
		attrsPrefix: h.attrsPrefix,
		groupPrefix: h.groupPrefix + name + ".",
		colors:      h.colors,
		ignoreKeys:  h.ignoreKeys,
		w:           h.w,
		mu:          h.mu,
	}
}

func (h *ColoredTextHandler) appendAttr(b *strings.Builder, attr slog.Attr, groupPrefix string) {
	if attr.Equal(slog.Attr{}) {
		return
	}
	if _, ok := h.ignoreKeys[groupPrefix+attr.Key]; ok {
		return
	}

	if attr.Value.Kind() == slog.KindGroup {
		if attr.Key != "" {
			groupPrefix += attr.Key + "."
		}
		for _, groupAttr := range attr.Value.Group() {
			h.appendAttr(b, groupAttr, groupPrefix)
		}
		return
	}
	if attr.Key == "" {
		return
	}

	b.WriteByte(' ')
	fmt.Fprintf(b, "%s%s=%s", h.colors.gray, groupPrefix+attr.Key, h.colors.reset)
	b.WriteString(formatValue(attr.Value.Resolve()))
}

func formatValue(v slog.Value) string {
	switch v.Kind() {
	case slog.KindAny:
		if data, err := json.Marshal(v.Any()); err == nil {
			return string(data)
		}
		return fmt.Sprintf("%v", v)
	case slog.KindBool:
		return strconv.FormatBool(v.Bool())
	case slog.KindDuration:
		return v.Duration().String()
	case slog.KindFloat64:
		return strconv.FormatFloat(v.Float64(), 'g', -1, 64)
	case slog.KindInt64:
		return strconv.FormatInt(v.Int64(), 10)
	case slog.KindString:
		return v.String()
	case slog.KindTime:
		return v.Time().Format(time.DateTime)
	case slog.KindUint64:
		return strconv.FormatUint(v.Uint64(), 10)
	default:
		return "unknown"
	}
}
