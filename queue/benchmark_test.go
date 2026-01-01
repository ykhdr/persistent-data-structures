package queue

import (
	"fmt"
	"testing"
)

var sinkInt int
var sinkBool bool

func BenchmarkQueueEnqueue(b *testing.B) {
	sizes := []int{100, 1000, 10000}

	for _, size := range sizes {
		b.Run(fmt.Sprintf("Queue/size_%d", size), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				q := NewQueue[int]()
				for j := 0; j < size; j++ {
					q = q.Enqueue(j)
				}
			}
		})

		b.Run(fmt.Sprintf("NaiveQueue/size_%d", size), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				q := NewNaiveQueue[int]()
				for j := 0; j < size; j++ {
					q = q.Enqueue(j)
				}
			}
		})
	}
}

func BenchmarkQueueDequeue(b *testing.B) {
	sizes := []int{100, 1000, 10000}

	for _, size := range sizes {
		b.Run(fmt.Sprintf("Queue/size_%d", size), func(b *testing.B) {
			q := NewQueue[int]()
			for j := 0; j < size; j++ {
				q = q.Enqueue(j)
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				var ok bool
				var v int
				q, v, ok = q.Dequeue()
				sinkInt = v
				sinkBool = ok
			}
		})

		b.Run(fmt.Sprintf("NaiveQueue/size_%d", size), func(b *testing.B) {
			q := NewNaiveQueue[int]()
			for j := 0; j < size; j++ {
				q = q.Enqueue(j)
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				var ok bool
				var v int
				q, v, ok = q.Dequeue()
				sinkInt = v
				sinkBool = ok
			}
		})
	}
}

func BenchmarkQueuePeek(b *testing.B) {
	sizes := []int{100, 1000, 10000, 100000}

	for _, size := range sizes {
		q1 := NewQueue[int]()
		q2 := NewNaiveQueue[int]()

		for i := 0; i < size; i++ {
			q1 = q1.Enqueue(i)
			q2 = q2.Enqueue(i)
		}

		b.Run(fmt.Sprintf("Queue/size_%d", size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				v, ok := q1.Peek()
				sinkInt = v
				sinkBool = ok
			}
		})

		b.Run(fmt.Sprintf("NaiveQueue/size_%d", size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				v, ok := q2.Peek()
				sinkInt = v
				sinkBool = ok
			}
		})
	}
}

func BenchmarkQueueVersionCreation(b *testing.B) {
	sizes := []int{1000, 10000}
	numVersions := []int{10, 100}

	for _, size := range sizes {
		baseQ := NewQueue[int]()
		baseN := NewNaiveQueue[int]()

		for i := 0; i < size; i++ {
			baseQ = baseQ.Enqueue(i)
			baseN = baseN.Enqueue(i)
		}

		for _, versions := range numVersions {
			b.Run(fmt.Sprintf("Queue/size_%d/versions_%d", size, versions), func(b *testing.B) {
				b.ReportAllocs()
				for i := 0; i < b.N; i++ {
					q := baseQ
					for j := 0; j < versions; j++ {
						q = q.Enqueue(j)
					}
				}
			})

			b.Run(fmt.Sprintf("NaiveQueue/size_%d/versions_%d", size, versions), func(b *testing.B) {
				b.ReportAllocs()
				for i := 0; i < b.N; i++ {
					q := baseN
					for j := 0; j < versions; j++ {
						q = q.Enqueue(j)
					}
				}
			})
		}
	}
}

func BenchmarkQueueMixedOperations(b *testing.B) {
	size := 1000

	b.Run("Queue", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			q := NewQueue[int]()

			for j := 0; j < size; j++ {
				q = q.Enqueue(j)
			}

			for j := 0; j < 100; j++ {
				_, _ = q.Peek()
				q = q.Enqueue(j)
			}

			for j := 0; j < 50; j++ {
				var v int
				q, v, _ = q.Dequeue()
				sinkInt = v
			}
		}
	})

	b.Run("NaiveQueue", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			q := NewNaiveQueue[int]()

			for j := 0; j < size; j++ {
				q = q.Enqueue(j)
			}

			for j := 0; j < 100; j++ {
				_, _ = q.Peek()
				q = q.Enqueue(j)
			}

			for j := 0; j < 50; j++ {
				var v int
				q, v, _ = q.Dequeue()
				sinkInt = v
			}
		}
	})
}

func BenchmarkQueueMemoryAllocation(b *testing.B) {
	size := 10000

	b.Run("Queue/Dequeue", func(b *testing.B) {
		q := NewQueue[int]()
		for i := 0; i < size; i++ {
			q = q.Enqueue(i)
		}

		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			var v int
			q, v, _ = q.Dequeue()
			sinkInt = v
		}
	})

	b.Run("NaiveQueue/Dequeue", func(b *testing.B) {
		q := NewNaiveQueue[int]()
		for i := 0; i < size; i++ {
			q = q.Enqueue(i)
		}

		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			var v int
			q, v, _ = q.Dequeue()
			sinkInt = v
		}
	})
}
