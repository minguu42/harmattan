package handler

import "github.com/minguu42/harmattan/internal/oapi"

func convertOptString(o oapi.OptString) *string {
	if o.Set {
		return &o.Value
	}
	return nil
}

func convertOptBool(o oapi.OptBool) *bool {
	if o.Set {
		return &o.Value
	}
	return nil
}
