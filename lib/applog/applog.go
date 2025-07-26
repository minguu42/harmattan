package applog

import (
	"context"
	"log/slog"
	"os"
	"reflect"
	"time"

	"github.com/minguu42/harmattan/api/apperr"
	"github.com/minguu42/harmattan/internal/auth"
	"github.com/minguu42/harmattan/lib/slogdebug"
)

type Logger struct {
	base *slog.Logger
}

func New(indented bool) *Logger {
	opts := &slog.HandlerOptions{ReplaceAttr: replaceAttr}
	if indented {
		return &Logger{base: slog.New(slogdebug.NewJSONIndentHandler(os.Stdout, opts))}
	}
	return &Logger{base: slog.New(slog.NewJSONHandler(os.Stdout, opts))}
}

func replaceAttr(_ []string, a slog.Attr) slog.Attr {
	if a.Key == slog.MessageKey {
		a.Key = "message"
	}

	if a.Value.Kind() != slog.KindAny {
		return a
	}
	switch rv := reflect.ValueOf(a.Value.Any()); rv.Kind() {
	case reflect.Pointer:
		if rv.IsNil() || rv.Elem().Kind() != reflect.Struct {
			return a
		}
		return slog.Any(a.Key, maskStructValue(rv.Elem()).Addr().Interface())
	case reflect.Struct:
		return slog.Any(a.Key, maskStructValue(rv).Interface())
	default:
		return a
	}
}

func maskStructValue(v reflect.Value) reflect.Value {
	t := v.Type()
	newStruct := reflect.New(t).Elem()
	for i := range v.NumField() {
		structField := t.Field(i)
		if !structField.IsExported() {
			continue
		}
		newField := newStruct.Field(i)
		if !newField.CanSet() {
			continue
		}

		if tag := structField.Tag.Get("log"); tag == "mask" {
			setMaskedValue(newField)
			continue
		}

		field := v.Field(i)
		switch {
		case field.Kind() == reflect.Struct:
			newField.Set(maskStructValue(field))
		case field.Kind() == reflect.Pointer && !field.IsNil() && field.Elem().Kind() == reflect.Struct:
			newField.Set(maskStructValue(field.Elem()).Addr())
		default:
			newField.Set(field)
		}
	}
	return newStruct
}

func setMaskedValue(newField reflect.Value) {
	t := newField.Type()
	switch t.Kind() {
	case reflect.Bool:
		newField.SetBool(false)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		newField.SetInt(0)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		newField.SetUint(0)
	case reflect.Float32, reflect.Float64:
		newField.SetFloat(0.0)
	case reflect.Map:
		newField.Set(reflect.MakeMap(t))
	case reflect.Slice:
		newField.Set(reflect.MakeSlice(t, 0, 0))
	case reflect.String:
		newField.SetString("<hidden>")
	default:
		newField.Set(reflect.Zero(t))
	}
}

type loggerKey struct{}

func (l *Logger) logger(ctx context.Context) *slog.Logger {
	if logger, ok := ctx.Value(loggerKey{}).(*slog.Logger); ok {
		return logger
	}
	return l.base
}

func (l *Logger) ContextWithRequestID(ctx context.Context, requestID string) context.Context {
	logger, ok := ctx.Value(loggerKey{}).(*slog.Logger)
	if !ok {
		logger = l.base
	}
	return context.WithValue(ctx, loggerKey{}, logger.With(slog.String("request_id", requestID)))
}

func (l *Logger) Event(ctx context.Context, msg string) {
	l.logger(ctx).Log(ctx, slog.LevelInfo, msg)
}

func (l *Logger) Error(ctx context.Context, msg string) {
	l.logger(ctx).Log(ctx, slog.LevelError, msg)
}

type AccessFields struct {
	ExecutionTime time.Duration
	Err           error

	OperationID string
	Method      string
	URL         string
	Body        any
	IPAddress   string
}

func (l *Logger) Access(ctx context.Context, fields *AccessFields) {
	message := "Request accepted"
	executionTime := slog.Int64("execution_time", fields.ExecutionTime.Milliseconds())
	var request slog.Attr
	if user, ok := auth.UserFromContext(ctx); ok {
		request = slog.Group("request",
			slog.String("user_id", string(user.ID)),
			slog.String("operation_id", fields.OperationID),
			slog.String("method", fields.Method),
			slog.String("url", fields.URL),
			slog.Any("body", fields.Body),
			slog.String("ip_address", fields.IPAddress),
		)
	} else {
		request = slog.Group("request",
			slog.String("operation_id", fields.OperationID),
			slog.String("method", fields.Method),
			slog.String("url", fields.URL),
			slog.Any("body", fields.Body),
			slog.String("ip_address", fields.IPAddress),
		)
	}

	if fields.Err == nil {
		l.logger(ctx).LogAttrs(ctx, slog.LevelInfo, message,
			slog.Int("status_code", 200),
			executionTime,
			request,
		)
		return
	}

	appErr := apperr.ToError(fields.Err)
	status := appErr.StatusCode()
	level := slog.LevelWarn
	if status >= 500 {
		level = slog.LevelError
	}
	if stacktrace := appErr.Stacktrace(); len(stacktrace) > 0 {
		l.logger(ctx).LogAttrs(ctx, level, message,
			slog.Int("status_code", status),
			slog.String("error_message", appErr.Error()),
			slog.Any("stacktrace", stacktrace),
			executionTime,
			request,
		)
		return
	}
	l.logger(ctx).LogAttrs(ctx, level, message,
		slog.Int("status_code", status),
		slog.String("error_message", appErr.Error()),
		executionTime,
		request,
	)
}
