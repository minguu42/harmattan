package usecase

type Option[T any] struct {
	V     T
	Valid bool
}
