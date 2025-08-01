package pointers

func Ref[T any](v T) *T {
	return &v
}

func OrZero[T any](v *T) T {
	if v != nil {
		return *v
	}
	var zero T
	return zero
}
