package array

import "iter"

type NaiveArray[T any] struct {
	data []T
}

func NewNaiveArray[T any]() *NaiveArray[T] {
	return &NaiveArray[T]{
		data: make([]T, 0),
	}
}

func (a *NaiveArray[T]) sealed() {}

func (a *NaiveArray[T]) Len() int {
	return len(a.data)
}

func (a *NaiveArray[T]) Get(index int) (T, bool) {
	var zero T
	if index < 0 || index >= len(a.data) {
		return zero, false
	}
	return a.data[index], true
}

func (a *NaiveArray[T]) Set(index int, value T) *NaiveArray[T] {
	if index < 0 || index >= len(a.data) {
		return a
	}

	newData := make([]T, len(a.data))
	copy(newData, a.data)
	newData[index] = value

	return &NaiveArray[T]{data: newData}
}

func (a *NaiveArray[T]) Append(value T) *NaiveArray[T] {
	newData := make([]T, len(a.data)+1)
	copy(newData, a.data)
	newData[len(a.data)] = value

	return &NaiveArray[T]{data: newData}
}

func (a *NaiveArray[T]) Pop() (*NaiveArray[T], T, bool) {
	var zero T
	if len(a.data) == 0 {
		return a, zero, false
	}

	value := a.data[len(a.data)-1]
	newData := make([]T, len(a.data)-1)
	copy(newData, a.data[:len(a.data)-1])

	return &NaiveArray[T]{data: newData}, value, true
}

func (a *NaiveArray[T]) All() iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		for i, v := range a.data {
			if !yield(i, v) {
				return
			}
		}
	}
}

func (a *NaiveArray[T]) Values() iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, v := range a.data {
			if !yield(v) {
				return
			}
		}
	}
}
