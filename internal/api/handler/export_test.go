package handler

var (
	ConvertOptDate     = convertOptDate
	ConvertOptDateTime = convertOptDateTime
)

func Ternary[T any](condition bool, trueVal, falseVal T) T {
	return ternary(condition, trueVal, falseVal)
}
