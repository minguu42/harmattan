package slogdebug

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"strings"
	"sync"
)

type JSONIndentHandler struct {
	handler slog.Handler
	w       io.Writer
	mu      *sync.Mutex
	buf     *bytes.Buffer
}

func NewJSONIndentHandler(w io.Writer, opts *slog.HandlerOptions) *JSONIndentHandler {
	buf := &bytes.Buffer{}
	return &JSONIndentHandler{
		handler: slog.NewJSONHandler(buf, opts),
		w:       w,
		mu:      &sync.Mutex{},
		buf:     buf,
	}
}

func (h *JSONIndentHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.handler.Enabled(ctx, level)
}

func (h *JSONIndentHandler) Handle(ctx context.Context, record slog.Record) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	if err := h.handler.Handle(ctx, record); err != nil {
		return err
	}

	encoder := json.NewEncoder(h.w)
	encoder.SetIndent("", strings.Repeat(" ", 2))
	if err := encoder.Encode(json.RawMessage(h.buf.Bytes())); err != nil {
		return fmt.Errorf("failed to encode log entry: %w", err)
	}
	h.buf.Reset()
	return nil
}

func (h *JSONIndentHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &JSONIndentHandler{
		handler: h.handler.WithAttrs(attrs),
		w:       h.w,
		mu:      h.mu,
		buf:     h.buf,
	}
}

func (h *JSONIndentHandler) WithGroup(name string) slog.Handler {
	return &JSONIndentHandler{
		handler: h.handler.WithGroup(name),
		w:       h.w,
		mu:      h.mu,
		buf:     h.buf,
	}
}
