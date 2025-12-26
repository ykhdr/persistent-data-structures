package queue

import "iter"

type stackNode[T any] struct {
	value T
	next  *stackNode[T]
}

type stack[T any] struct {
	head *stackNode[T]
	len  int
}

func newStack[T any]() *stack[T] {
	return &stack[T]{}
}

func (s *stack[T]) push(value T) *stack[T] {
	return &stack[T]{
		head: &stackNode[T]{value: value, next: s.head},
		len:  s.len + 1,
	}
}

func (s *stack[T]) pop() (*stack[T], T, bool) {
	var zero T
	if s.head == nil {
		return s, zero, false
	}
	return &stack[T]{
		head: s.head.next,
		len:  s.len - 1,
	}, s.head.value, true
}

func (s *stack[T]) peek() (T, bool) {
	var zero T
	if s.head == nil {
		return zero, false
	}
	return s.head.value, true
}

func (s *stack[T]) isEmpty() bool {
	return s.head == nil
}

func (s *stack[T]) reverse() *stack[T] {
	result := newStack[T]()
	current := s
	for !current.isEmpty() {
		var value T
		current, value, _ = current.pop()
		result = result.push(value)
	}
	return result
}

type Queue[T any] struct {
	front *stack[T]
	rear  *stack[T]
	len   int
}

func NewQueue[T any]() *Queue[T] {
	return &Queue[T]{
		front: newStack[T](),
		rear:  newStack[T](),
		len:   0,
	}
}

func (q *Queue[T]) Len() int {
	return q.len
}

func (q *Queue[T]) IsEmpty() bool {
	return q.len == 0
}

func (q *Queue[T]) Enqueue(value T) *Queue[T] {
	return &Queue[T]{
		front: q.front,
		rear:  q.rear.push(value),
		len:   q.len + 1,
	}
}

func (q *Queue[T]) Dequeue() (*Queue[T], T, bool) {
	var zero T
	if q.len == 0 {
		return q, zero, false
	}

	front := q.front
	rear := q.rear

	if front.isEmpty() {
		front = rear.reverse()
		rear = newStack[T]()
	}

	newFront, value, _ := front.pop()

	return &Queue[T]{
		front: newFront,
		rear:  rear,
		len:   q.len - 1,
	}, value, true
}

func (q *Queue[T]) Peek() (T, bool) {
	var zero T
	if q.len == 0 {
		return zero, false
	}

	front := q.front
	if front.isEmpty() {
		front = q.rear.reverse()
	}

	return front.peek()
}

func (q *Queue[T]) All() iter.Seq[T] {
	return func(yield func(T) bool) {
		current := q
		for current.len > 0 {
			var value T
			current, value, _ = current.Dequeue()
			if !yield(value) {
				return
			}
		}
	}
}
