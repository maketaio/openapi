package ptr

func Deref[T any](v *T, def T) T {
	if v == nil {
		return def
	}

	return *v
}

func To[T any](v T) *T {
	return &v
}
