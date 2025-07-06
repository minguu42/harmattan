package handler

import "github.com/minguu42/harmattan/internal/openapi"

func convertOptString(o openapi.OptString) *string {
	if o.Set {
		return &o.Value
	}
	return nil
}

func convertOptBool(o openapi.OptBool) *bool {
	if o.Set {
		return &o.Value
	}
	return nil
}
