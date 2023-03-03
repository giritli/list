package list

import (
	"sort"
)

type Of[T any] []T

func Convert[From, To any](from Of[From], fn func(From) To) Of[To] {
	to := make(Of[To], 0, len(from))
	for _, item := range from {
		to = append(to, fn(item))
	}

	return to
}

func ReduceInto[From, To any](from Of[From], fn func(To, From) To) To {
	var zero To
	if len(from) == 0 {
		return zero
	}

	zero = fn(zero, from[0])
	for i, j := 1, len(from); i < j; i++ {
		zero = fn(zero, from[i])
	}

	return zero
}

func (l Of[T]) Contains(item T, fn func(T, T) bool) bool {
	for _, v := range l {
		if fn(v, item) {
			return true
		}
	}

	return false
}

func (l Of[T]) Chunk(size uint) []Of[T] {
	if size == 0 {
		return []Of[T]{l}
	}

	chunks := len(l) / int(size)
	remainder := len(l) % int(size)
	total := chunks
	if remainder > 0 {
		total++
	}

	nls := make([]Of[T], total)
	for i, item := range l {
		nls[i/int(size)] = append(nls[i/int(size)], item)
	}

	return nls
}

func (l Of[T]) Filter(fn func(T) bool) Of[T] {
	nl := Of[T]{}
	for _, item := range l {
		if fn(item) {
			nl = append(nl, item)
		}
	}

	return nl
}

func (l Of[T]) Insert(at uint, items ...T) Of[T] {
	if len(items) == 0 {
		return l
	}

	nl := make(Of[T], 0, len(l)+len(items))
	if at == 0 {
		return append(append(nl, items...), l...)
	}

	if int(at) >= len(l) {
		return append(append(nl, l...), items...)
	}

	return append(append(append(nl, l[:at]...), items...), l[at:]...)
}

func (l Of[T]) Map(fn func(T) T) Of[T] {
	nl := make(Of[T], 0, len(l))
	for _, item := range l {
		nl = append(nl, fn(item))
	}

	return nl
}

func (l Of[T]) Sort(fn func(a, b T) bool) Of[T] {
	nl := l
	sort.SliceStable(nl, func(i, j int) bool {
		return fn(l[i], l[j])
	})

	return nl
}

/*
func (l Of[T]) Unique() Of[T] {
	if len(l) <= 1 {
		return l
	}

	valueMap := map[T]int{}
	ri := 0
	for _, v := range l {
		if _, ok := valueMap[v]; !ok {
			valueMap[v] = ri
		} else {
			ri--
		}

		ri++
	}

	nl := make(Of[T], len(valueMap))
	for k, v := range valueMap {
		nl[v] = k
	}

	return nl
}
*/
