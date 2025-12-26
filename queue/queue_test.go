package queue

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/ykhdr/persistent-data-structures/array"
)

func TestQueue_EnqueueDequeue(t *testing.T) {
	t.Run("добавление и извлечение элементов", func(t *testing.T) {
		q := NewQueue[int]()

		q = q.Enqueue(1)
		q = q.Enqueue(2)
		q = q.Enqueue(3)

		assert.Equal(t, 3, q.Len(), "после добавления 3 элементов длина должна быть 3")

		q, val, ok := q.Dequeue()
		require.True(t, ok, "Dequeue должен вернуть ok")
		assert.Equal(t, 1, val, "первый извлечённый элемент должен быть 1")

		q, val, ok = q.Dequeue()
		require.True(t, ok)
		assert.Equal(t, 2, val, "второй извлечённый элемент должен быть 2")

		q, val, ok = q.Dequeue()
		require.True(t, ok)
		assert.Equal(t, 3, val, "третий извлечённый элемент должен быть 3")

		assert.Equal(t, 0, q.Len(), "после извлечения всех элементов очередь должна быть пустой")
	})
}

func TestQueue_DequeueEmpty(t *testing.T) {
	t.Run("извлечение из пустой очереди", func(t *testing.T) {
		q := NewQueue[int]()

		_, _, ok := q.Dequeue()

		assert.False(t, ok, "Dequeue из пустой очереди должен вернуть false")
	})
}

func TestQueue_Persistence(t *testing.T) {
	t.Run("Enqueue создаёт новую версию", func(t *testing.T) {
		q1 := NewQueue[int]()
		q2 := q1.Enqueue(1)
		q3 := q2.Enqueue(2)

		assert.Equal(t, 0, q1.Len(), "q1 не должен измениться")
		assert.Equal(t, 1, q2.Len(), "q2 не должен измениться")
		assert.Equal(t, 2, q3.Len(), "q3 должен содержать 2 элемента")
	})

	t.Run("Dequeue создаёт новую версию", func(t *testing.T) {
		q1 := NewQueue[int]().Enqueue(1).Enqueue(2)

		q2, val, _ := q1.Dequeue()

		assert.Equal(t, 1, val)
		assert.Equal(t, 2, q1.Len(), "оригинал не должен измениться после Dequeue")
		assert.Equal(t, 1, q2.Len(), "новая версия должна содержать 1 элемент")
	})
}

func TestQueue_Peek(t *testing.T) {
	t.Run("просмотр первого элемента", func(t *testing.T) {
		q := NewQueue[int]().Enqueue(1).Enqueue(2)

		val, ok := q.Peek()

		require.True(t, ok, "Peek должен вернуть ok")
		assert.Equal(t, 1, val, "Peek должен вернуть первый элемент")
		assert.Equal(t, 2, q.Len(), "Peek не должен изменять очередь")
	})

	t.Run("просмотр пустой очереди", func(t *testing.T) {
		q := NewQueue[int]()

		_, ok := q.Peek()

		assert.False(t, ok, "Peek пустой очереди должен вернуть false")
	})
}

func TestQueue_IsEmpty(t *testing.T) {
	t.Run("новая очередь пуста", func(t *testing.T) {
		q := NewQueue[int]()

		assert.True(t, q.IsEmpty(), "новая очередь должна быть пустой")
	})

	t.Run("очередь с элементами не пуста", func(t *testing.T) {
		q := NewQueue[int]().Enqueue(1)

		assert.False(t, q.IsEmpty(), "очередь с элементом не должна быть пустой")
	})

	t.Run("после извлечения всех элементов очередь пуста", func(t *testing.T) {
		q := NewQueue[int]().Enqueue(1)
		q, _, _ = q.Dequeue()

		assert.True(t, q.IsEmpty(), "после извлечения всех элементов очередь должна быть пустой")
	})
}

func TestQueue_FIFOOrder(t *testing.T) {
	t.Run("порядок FIFO сохраняется", func(t *testing.T) {
		q := NewQueue[int]()

		for i := 1; i <= 100; i++ {
			q = q.Enqueue(i)
		}

		for i := 1; i <= 100; i++ {
			var val int
			q, val, _ = q.Dequeue()
			assert.Equal(t, i, val, "элемент %d должен быть извлечён в правильном порядке", i)
		}
	})
}

func TestQueue_MixedOperations(t *testing.T) {
	t.Run("чередование добавления и извлечения", func(t *testing.T) {
		q := NewQueue[int]()

		q = q.Enqueue(1)
		q = q.Enqueue(2)

		q, val, _ := q.Dequeue()
		assert.Equal(t, 1, val, "первое извлечение должно вернуть 1")

		q = q.Enqueue(3)
		q = q.Enqueue(4)

		q, val, _ = q.Dequeue()
		assert.Equal(t, 2, val, "второе извлечение должно вернуть 2")

		q, val, _ = q.Dequeue()
		assert.Equal(t, 3, val, "третье извлечение должно вернуть 3")

		q, val, _ = q.Dequeue()
		assert.Equal(t, 4, val, "четвёртое извлечение должно вернуть 4")
	})
}

func TestQueue_Iterator(t *testing.T) {
	t.Run("итерация по всем элементам", func(t *testing.T) {
		q := NewQueue[int]().Enqueue(1).Enqueue(2).Enqueue(3)

		sum := 0
		count := 0
		for val := range q.All() {
			sum += val
			count++
		}

		assert.Equal(t, 3, count, "должно быть 3 элемента")
		assert.Equal(t, 6, sum, "сумма должна быть 6")
	})
}

func TestQueue_NestedStructures(t *testing.T) {
	t.Run("очередь векторов", func(t *testing.T) {
		inner := array.NewVector[int]().Append(1).Append(2)
		q := NewQueue[*array.Vector[int]]().Enqueue(inner)

		_, retrieved, ok := q.Dequeue()

		require.True(t, ok, "Dequeue должен вернуть ok")

		val, _ := retrieved.Get(0)
		assert.Equal(t, 1, val, "вложенный вектор должен содержать правильные данные")
	})
}
