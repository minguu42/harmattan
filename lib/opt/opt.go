package opt

type Option[T any] struct {
	V     T
	Valid bool
}

func FromPointer[T any](p *T) Option[T] {
	if p == nil {
		return Option[T]{Valid: false}
	}
	return Option[T]{V: *p, Valid: true}
}

func (o Option[T]) ToPointer() *T {
	if o.Valid {
		return &o.V
	}
	return nil
}
