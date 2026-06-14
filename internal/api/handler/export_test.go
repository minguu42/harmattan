package handler

var (
	ConvertOptDate      = convertOptDate
	ConvertOptDateTime  = convertOptDateTime
	ValidateEmail       = validateEmail
	ValidatePassword    = validatePassword
	ValidateProjectName = validateProjectName
	ValidateTaskName    = validateTaskName
	ValidateStepName    = validateStepName
	ValidateTagName     = validateTagName
)

func Ternary[T any](condition bool, trueVal, falseVal T) T {
	return ternary(condition, trueVal, falseVal)
}
