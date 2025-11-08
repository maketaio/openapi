package fields

type Optional[T any] struct{ s State[T] }

func OptionalValue[T any](v T) Optional[T] {
	return Optional[T]{
		s: State[T]{v: &v, isSet: true},
	}
}

func OptionalUnset[T any]() Optional[T] {
	return Optional[T]{
		s: State[T]{isSet: false},
	}
}

func (o Optional[T]) State() State[T] {
	return o.s
}

func (o *Optional[T]) Set(v T) {
	o.s.Set(v)
}

func (o *Optional[T]) Unset() {
	o.s.Unset()
}

func (o Optional[T]) IsPresent() bool {
	return o.s.IsPresent()
}

func (o Optional[T]) Value() (T, bool) {
	return o.s.Value()
}

func (o Optional[T]) IsZero() bool {
	return !o.s.IsPresent()
}

func (o *Optional[T]) UnmarshalJSON(data []byte) error {
	return o.s.UnmarshalJSON(data)
}

func (o Optional[T]) MarshalJSON() ([]byte, error) {
	return o.s.MarshalJSON()
}
