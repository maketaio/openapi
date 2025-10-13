package set

type empty struct{}

type Set[T comparable] map[T]empty

func NewSet[T comparable]() Set[T] {
	return Set[T]{}
}

func (s Set[T]) Has(item T) bool {
	_, exists := s[item]
	return exists
}

func (s Set[T]) Add(item T) {
	s[item] = empty{}
}

func (s Set[T]) Delete(item T) {
	delete(s, item)
}

func (s Set[T]) Len() int {
	return len(s)
}

func (s Set[T]) Merge(other Set[T]) {
	for k := range other {
		s.Add(k)
	}
}
