package handler_test

import (
	"testing"

	"github.com/minguu42/harmattan/internal/api/openapi"
	"github.com/minguu42/harmattan/internal/lib/httpcheck"
)

func TestHandler_CheckHealth(t *testing.T) {
	want := openapi.CheckHealthOK{Revision: "xxxxxxx"}
	httpcheck.New(th).Test(t, "GET", "/health").
		Check().HasStatus(200).HasJSON(want)
}
