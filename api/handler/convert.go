package handler

import (
	"time"

	"github.com/minguu42/harmattan/internal/openapi"
)

func convertOptString(o openapi.OptString) *string {
	if o.Set {
		return &o.Value
	}
	return nil
}

func convertTimePtr(v *time.Time) openapi.OptDateTime {
	if v != nil {
		return openapi.OptDateTime{Value: *v, Set: true}
	}
	return openapi.OptDateTime{}
}

func convertOptBool(o openapi.OptBool) *bool {
	if o.Set {
		return &o.Value
	}
	return nil
}
