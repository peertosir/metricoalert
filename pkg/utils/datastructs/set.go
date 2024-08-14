package datastructs

type HashSet[T comparable] map[T]struct{}

func NewHashSet[T comparable]() HashSet[T] {
	return HashSet[T]{}
}

func NewHashSetWithValues[T comparable](values ...T) HashSet[T] {
	hs := HashSet[T]{}
	for _, v := range values {
		hs.Add(v)
	}
	return hs
}

func (hs HashSet[T]) Add(el T) {
	hs[el] = struct{}{}
}

func (hs HashSet[T]) Contains(el T) bool {
	_, ok := hs[el]
	return ok
}
