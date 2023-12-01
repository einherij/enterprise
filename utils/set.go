package utils

type Set[data comparable] map[data]struct{}

func NewSet[T comparable]() Set[T] {
	return Set[T]{}
}

func (s Set[T]) Add(data ...T) {
	for _, d := range data {
		s[d] = struct{}{}
	}
}

func (s Set[T]) Delete(data T) {
	delete(s, data)
}

func (s Set[T]) Exists(data T) bool {
	_, ok := s[data]
	return ok
}

func (s Set[T]) Values() []T {
	var values []T
	for value := range s {
		values = append(values, value)
	}
	return values
}
