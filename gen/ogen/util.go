package ogen

import "time"

func (d OptDate) Ptr() *time.Time {
	if d.IsSet() {
		return &d.Value
	}
	return nil
}
