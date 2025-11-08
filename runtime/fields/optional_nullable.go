package fields

type OptionalNullable[T any] struct {
	s State[T]
}

func OptionalNullableValue[T any](v T) OptionalNullable[T] {
	return OptionalNullable[T]{
		s: State[T]{v: &v, isSet: true},
	}
}

func OptionalNullableUnset[T any]() OptionalNullable[T] {
	return OptionalNullable[T]{
		s: State[T]{isSet: false},
	}
}

func OptionalNull[T any]() OptionalNullable[T] {
	return OptionalNullable[T]{
		s: State[T]{v: nil, isSet: true},
	}
}

func (o OptionalNullable[T]) State() State[T] {
	return o.s
}

func (o *OptionalNullable[T]) Set(v T) {
	o.s.Set(v)
}

func (o *OptionalNullable[T]) SetNull() {
	o.s.SetNull()
}

func (o *OptionalNullable[T]) Unset() {
	o.s.Unset()
}

func (o OptionalNullable[T]) IsPresent() bool {
	return o.s.IsPresent()
}

func (o OptionalNullable[T]) IsNull() bool {
	return o.s.IsNull()
}

func (o OptionalNullable[T]) Value() (T, bool) {
	return o.s.Value()
}

func (o OptionalNullable[T]) IsZero() bool {
	return !o.s.IsPresent()
}

func (o *OptionalNullable[T]) UnmarshalJSON(data []byte) error {
	return o.s.UnmarshalJSON(data)
}

func (o OptionalNullable[T]) MarshalJSON() ([]byte, error) {
	return o.s.MarshalJSON()
}
