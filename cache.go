package ganrac

type Cacher[T any] interface {
	Put(key Hashable, value T)
	Get(key Hashable) (T, bool)
	Len() int
}

type NoCache[T any] struct {
	cnt int
}

func (nc NoCache[T]) Put(key Hashable, value T) {
}

func (nc NoCache[T]) Get(key Hashable) (T, bool) {
	nc.cnt++
	var zero T
	return zero, false
}

func (nc NoCache[T]) Len() int {
	return 0
}

func (nc NoCache[T]) GetCount() int {
	return nc.cnt
}
