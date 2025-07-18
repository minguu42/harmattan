package handler_test

import (
	"testing"

	"github.com/ikawaha/httpcheck"
	"github.com/minguu42/harmattan/internal/openapi"
)

func TestHandler_CheckHealth(t *testing.T) {
	want := openapi.CheckHealthOK{Revision: "xxxxxxx"}
	httpcheck.New(th).Test(t, "GET", "/health").
		Check().HasStatus(200).HasJSON(want)
}
