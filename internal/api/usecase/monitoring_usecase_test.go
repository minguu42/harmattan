package usecase_test

import (
	"testing"

	"github.com/minguu42/harmattan/internal/api/openapi"
)

func TestMonitoring_CheckHealth(t *testing.T) {
	runTest(t, test{
		Method:     "GET",
		Path:       "/health",
		WantStatus: 200,
		WantJSON:   openapi.CheckHealthOK{Revision: testRevision},
	})
}
