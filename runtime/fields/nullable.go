package fields

type Nullable[T any] struct{ s State[T] }

func NullableValue[T any](v T) Nullable[T] {
	return Nullable[T]{
		s: State[T]{v: &v, isSet: true},
	}
}

func Null[T any]() Nullable[T] {
	return Nullable[T]{
		s: State[T]{v: nil, isSet: true},
	}
}

func (r Nullable[T]) State() State[T] {
	return r.s
}

func (r *Nullable[T]) Set(v T) {
	r.s.Set(v)
}

func (r *Nullable[T]) SetNull() {
	r.s.SetNull()
}

func (r Nullable[T]) IsNull() bool {
	return r.s.IsNull()
}

func (r Nullable[T]) Value() (T, bool) {
	return r.s.Value()
}

func (r Nullable[T]) IsZero() bool {
	return false
}

func (r *Nullable[T]) UnmarshalJSON(data []byte) error {
	return r.s.UnmarshalJSON(data)
}

func (r Nullable[T]) MarshalJSON() ([]byte, error) {
	return r.s.MarshalJSON()
}
