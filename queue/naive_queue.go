package queue

type NaiveQueue[T any] struct {
	data []T
}

func NewNaiveQueue[T any]() *NaiveQueue[T] {
	return &NaiveQueue[T]{data: make([]T, 0)}
}

func (q *NaiveQueue[T]) Len() int { return len(q.data) }

func (q *NaiveQueue[T]) IsEmpty() bool { return len(q.data) == 0 }

func (q *NaiveQueue[T]) Enqueue(value T) *NaiveQueue[T] {
	newData := make([]T, len(q.data)+1)
	copy(newData, q.data)
	newData[len(q.data)] = value
	return &NaiveQueue[T]{data: newData}
}

func (q *NaiveQueue[T]) Dequeue() (*NaiveQueue[T], T, bool) {
	var zero T
	if len(q.data) == 0 {
		return q, zero, false
	}

	value := q.data[0]
	newData := make([]T, len(q.data)-1)
	copy(newData, q.data[1:])

	return &NaiveQueue[T]{data: newData}, value, true
}

func (q *NaiveQueue[T]) Peek() (T, bool) {
	var zero T
	if len(q.data) == 0 {
		return zero, false
	}
	return q.data[0], true
}
