package handler

import (
	"time"

	"github.com/minguu42/harmattan/internal/api/openapi"
	"github.com/minguu42/harmattan/internal/api/usecase"
)

type Handler struct {
	openapi.UnimplementedHandler
	Authentication usecase.Authentication
	Monitoring     usecase.Monitoring
	Project        usecase.Project
	Step           usecase.Step
	Tag            usecase.Tag
	Task           usecase.Task
}

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func ternary[T any](condition bool, trueVal, falseVal T) T {
	if condition {
		return trueVal
	}
	return falseVal
}

func convertOptDate(t *time.Time) openapi.OptDate {
	if t != nil {
		return openapi.OptDate{Value: *t, Set: true}
	}
	return openapi.OptDate{}
}

func convertOptDateTime(t *time.Time) openapi.OptDateTime {
	if t != nil {
		return openapi.OptDateTime{Value: *t, Set: true}
	}
	return openapi.OptDateTime{}
}

func convertSlice[T ~string](s []string) []T {
	r := make([]T, 0, len(s))
	for _, e := range s {
		r = append(r, T(e))
	}
	return r
}
