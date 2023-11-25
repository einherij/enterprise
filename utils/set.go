package utils

type Set[data comparable] map[data]struct{}

func NewSet[dataType comparable]() Set[dataType] {
	return Set[dataType]{}
}

func (s Set[dataType]) Add(data ...dataType) {
	for _, d := range data {
		s[d] = struct{}{}
	}
}

func (s Set[dataType]) Delete(data dataType) {
	delete(s, data)
}

func (s Set[dataType]) Exists(data dataType) bool {
	_, ok := s[data]
	return ok
}

func (s Set[dataType]) Values() []dataType {
	var values []dataType
	for value := range s {
		values = append(values, value)
	}
	return values
}
