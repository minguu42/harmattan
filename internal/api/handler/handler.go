package handler

import (
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
