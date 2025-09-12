package pointers

func Ref[T any](v T) *T {
	return &v
}

func Ternary[T any](condition bool, trueVal, falseVal T) T {
	if condition {
		return trueVal
	}
	return falseVal
}

func OrZero[T any](v *T) T {
	if v != nil {
		return *v
	}
	var zero T
	return zero
}
