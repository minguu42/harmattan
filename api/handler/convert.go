package handler

import (
	"time"

	"github.com/minguu42/harmattan/internal/domain"
	"github.com/minguu42/harmattan/internal/openapi"
	"github.com/minguu42/harmattan/lib/pointers"
)

func convertOptBool(v openapi.OptBool) *bool {
	if v.Set {
		return &v.Value
	}
	return nil
}

func convertOptInt(o openapi.OptInt) *int {
	if o.Set {
		return &o.Value
	}
	return nil
}

func convertOptString(v openapi.OptString) *string {
	if v.Set {
		return &v.Value
	}
	return nil
}

func convertOptColorString(v openapi.OptUpdateProjectReqColor) *domain.ProjectColor {
	if v.Set {
		return pointers.Ref(domain.ProjectColor(v.Value))
	}
	return nil
}

func convertDatePtr(v *time.Time) openapi.OptDate {
	if v != nil {
		return openapi.OptDate{Value: *v, Set: true}
	}
	return openapi.OptDate{}
}

func convertOptDateTime(v openapi.OptDateTime) *time.Time {
	if v.Set {
		return &v.Value
	}
	return nil
}

func convertDateTimePtr(v *time.Time) openapi.OptDateTime {
	if v != nil {
		return openapi.OptDateTime{Value: *v, Set: true}
	}
	return openapi.OptDateTime{}
}
