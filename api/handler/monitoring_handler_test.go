package handler_test

import (
	"testing"

	"github.com/ikawaha/httpcheck"
	"github.com/minguu42/harmattan/internal/oapi"
)

func TestHandler_CheckHealth(t *testing.T) {
	wantResponse := &oapi.CheckHealthOK{Revision: "xxxxxxx"}

	checker := httpcheck.New(h)
	checker.Test(t, "GET", "/health").
		Check().
		HasStatus(200).
		HasJSON(wantResponse)
}
