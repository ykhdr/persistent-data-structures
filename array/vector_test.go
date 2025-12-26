package array

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVector_NewVector(t *testing.T) {
	v := NewVector[int]()

	assert.NotNil(t, v)
	assert.Equal(t, 0, v.Len(), "новый вектор должен быть пустым")
}

func TestVector_Append(t *testing.T) {
	t.Run("добавление одного элемента", func(t *testing.T) {
		v := NewVector[int]().Append(42)

		assert.Equal(t, 1, v.Len())
		val, ok := v.Get(0)
		require.True(t, ok, "элемент должен существовать")
		assert.Equal(t, 42, val)
	})

	t.Run("добавление множества элементов", func(t *testing.T) {
		v := NewVector[int]()
		for i := 0; i < 100; i++ {
			v = v.Append(i)
		}

		assert.Equal(t, 100, v.Len())
		for i := 0; i < 100; i++ {
			val, ok := v.Get(i)
			require.True(t, ok, "Get(%d) должен вернуть ok", i)
			assert.Equal(t, i, val, "Get(%d) должен вернуть %d", i, i)
		}
	})

	t.Run("цепочка вызовов", func(t *testing.T) {
		v := NewVector[string]().
			Append("a").
			Append("b").
			Append("c")

		assert.Equal(t, 3, v.Len())

		val, _ := v.Get(0)
		assert.Equal(t, "a", val)

		val, _ = v.Get(1)
		assert.Equal(t, "b", val)

		val, _ = v.Get(2)
		assert.Equal(t, "c", val)
	})
}

func TestVector_Get(t *testing.T) {
	v := NewVector[int]().Append(10).Append(20).Append(30)

	t.Run("корректные индексы", func(t *testing.T) {
		val, ok := v.Get(0)
		assert.True(t, ok)
		assert.Equal(t, 10, val)

		val, ok = v.Get(1)
		assert.True(t, ok)
		assert.Equal(t, 20, val)

		val, ok = v.Get(2)
		assert.True(t, ok)
		assert.Equal(t, 30, val)
	})

	t.Run("отрицательный индекс", func(t *testing.T) {
		_, ok := v.Get(-1)
		assert.False(t, ok, "отрицательный индекс должен вернуть false")
	})

	t.Run("индекс за пределами", func(t *testing.T) {
		_, ok := v.Get(3)
		assert.False(t, ok, "индекс за пределами должен вернуть false")

		_, ok = v.Get(100)
		assert.False(t, ok, "индекс далеко за пределами - должен вернуть false")
	})

	t.Run("пустой вектор", func(t *testing.T) {
		empty := NewVector[int]()
		_, ok := empty.Get(0)
		assert.False(t, ok, "Get на пустом векторе должен вернуть false")
	})
}

func TestVector_Set(t *testing.T) {
	v := NewVector[int]().Append(1).Append(2).Append(3)

	t.Run("корректное изменение", func(t *testing.T) {
		v2 := v.Set(1, 100)

		val, _ := v2.Get(1)
		assert.Equal(t, 100, val)

		originalVal, _ := v.Get(1)
		assert.Equal(t, 2, originalVal, "оригинал не должен измениться")
	})

	t.Run("изменение первого элемента", func(t *testing.T) {
		v2 := v.Set(0, 999)

		val, _ := v2.Get(0)
		assert.Equal(t, 999, val)
	})

	t.Run("изменение последнего элемента", func(t *testing.T) {
		v2 := v.Set(2, 888)

		val, _ := v2.Get(2)
		assert.Equal(t, 888, val)
	})

	t.Run("некорректный индекс возвращает тот же вектор", func(t *testing.T) {
		v2 := v.Set(-1, 100)
		assert.Equal(t, v, v2, "Set с отрицательным индексом должен вернуть тот же вектор")

		v3 := v.Set(100, 100)
		assert.Equal(t, v, v3, "Set с индексом за пределами должен вернуть тот же вектор")
	})
}

func TestVector_Pop(t *testing.T) {
	t.Run("удаление из непустого вектора", func(t *testing.T) {
		v := NewVector[int]().Append(1).Append(2).Append(3)

		v2, val, ok := v.Pop()

		require.True(t, ok, "Pop должен вернуть ok")
		assert.Equal(t, 3, val, "Pop должен вернуть последний элемент")
		assert.Equal(t, 2, v2.Len(), "длина после Pop должна уменьшиться")
		assert.Equal(t, 3, v.Len(), "оригинал не должен измениться")
	})

	t.Run("удаление всех элементов", func(t *testing.T) {
		v := NewVector[int]().Append(1).Append(2)

		v, val, ok := v.Pop()
		require.True(t, ok)
		assert.Equal(t, 2, val)

		v, val, ok = v.Pop()
		require.True(t, ok)
		assert.Equal(t, 1, val)

		assert.Equal(t, 0, v.Len(), "после удаления всех элементов длина должна быть 0")
	})

	t.Run("удаление из пустого вектора", func(t *testing.T) {
		v := NewVector[int]()

		v2, _, ok := v.Pop()

		assert.False(t, ok, "Pop из пустого вектора должен вернуть false")
		assert.Equal(t, v, v2, "Pop из пустого вектора должен вернуть тот же вектор")
	})

	t.Run("удаление единственного элемента", func(t *testing.T) {
		v := NewVector[int]().Append(42)

		v2, val, ok := v.Pop()

		require.True(t, ok)
		assert.Equal(t, 42, val)
		assert.Equal(t, 0, v2.Len())
	})
}

func TestVector_Persistence(t *testing.T) {
	t.Run("Append создаёт новую копию", func(t *testing.T) {
		v1 := NewVector[int]()
		v2 := v1.Append(1)
		v3 := v2.Append(2)

		assert.Equal(t, 0, v1.Len(), "v1 не должен измениться")
		assert.Equal(t, 1, v2.Len(), "v2 не должен измениться")
		assert.Equal(t, 2, v3.Len(), "v3 должен содержать 2 элемента")
	})

	t.Run("Set создаёт новую копию", func(t *testing.T) {
		v1 := NewVector[int]().Append(1).Append(2).Append(3)
		v2 := v1.Set(1, 100)

		val1, _ := v1.Get(1)
		val2, _ := v2.Get(1)

		assert.Equal(t, 2, val1, "оригинал не должен измениться")
		assert.Equal(t, 100, val2, "новая версия должна содержать изменение")
	})

	t.Run("Pop создаёт новую копию", func(t *testing.T) {
		v1 := NewVector[int]().Append(1).Append(2)
		v2, _, _ := v1.Pop()

		assert.Equal(t, 2, v1.Len(), "оригинал не должен измениться")
		assert.Equal(t, 1, v2.Len(), "новая версия должна быть короче")
	})

	t.Run("ветвление версий", func(t *testing.T) {
		base := NewVector[int]().Append(1).Append(2)

		branch1 := base.Append(3)
		branch2 := base.Append(100)

		assert.Equal(t, 3, branch1.Len())
		assert.Equal(t, 3, branch2.Len())

		val1, _ := branch1.Get(2)
		val2, _ := branch2.Get(2)

		assert.Equal(t, 3, val1, "ветка 1 должна содержать 3")
		assert.Equal(t, 100, val2, "ветка 2 должна содержать 100")
	})
}

func TestVector_LargeDataset(t *testing.T) {
	sizes := []int{100, 1000, 10000}

	for _, size := range sizes {
		t.Run(fmt.Sprintf("size_%d", size), func(t *testing.T) {
			v := NewVector[int]()

			for i := 0; i < size; i++ {
				v = v.Append(i)
			}

			assert.Equal(t, size, v.Len())

			for i := 0; i < size; i++ {
				val, ok := v.Get(i)
				require.True(t, ok, "элемент %d должен существовать", i)
				assert.Equal(t, i, val, "элемент %d должен быть равен %d", i, i)
			}
		})
	}
}

func TestVector_Iterator(t *testing.T) {
	t.Run("итератор All", func(t *testing.T) {
		v := NewVector[int]().Append(10).Append(20).Append(30)

		var indices []int
		var values []int

		for i, val := range v.All() {
			indices = append(indices, i)
			values = append(values, val)
		}

		assert.Equal(t, []int{0, 1, 2}, indices, "индексы должны быть 0, 1, 2")
		assert.Equal(t, []int{10, 20, 30}, values, "значения должны быть 10, 20, 30")
	})

	t.Run("итератор Values", func(t *testing.T) {
		v := NewVector[int]().Append(1).Append(2).Append(3)

		sum := 0
		for val := range v.Values() {
			sum += val
		}

		assert.Equal(t, 6, sum, "сумма должна быть 6")
	})

	t.Run("итерация пустого вектора", func(t *testing.T) {
		v := NewVector[int]()

		count := 0
		for range v.All() {
			count++
		}

		assert.Equal(t, 0, count, "итерация пустого вектора не должна выполняться")
	})

	t.Run("ранний выход из итератора", func(t *testing.T) {
		v := NewVector[int]().Append(1).Append(2).Append(3).Append(4).Append(5)

		count := 0
		for _, val := range v.All() {
			count++
			if val == 3 {
				break
			}
		}

		assert.Equal(t, 3, count, "должно быть обработано 3 элемента до break")
	})
}

func TestVector_NestedStructures(t *testing.T) {
	t.Run("вектор векторов", func(t *testing.T) {
		inner1 := NewVector[int]().Append(1).Append(2)
		inner2 := NewVector[int]().Append(3).Append(4)

		outer := NewVector[*Vector[int]]().Append(inner1).Append(inner2)

		assert.Equal(t, 2, outer.Len())

		retrieved, ok := outer.Get(0)
		require.True(t, ok, "внутренний вектор должен существовать")
		assert.Equal(t, 2, retrieved.Len())

		val, _ := retrieved.Get(0)
		assert.Equal(t, 1, val)
	})

	t.Run("вектор структур", func(t *testing.T) {
		type Person struct {
			Name string
			Age  int
		}

		v := NewVector[Person]().
			Append(Person{"Alice", 30}).
			Append(Person{"Bob", 25})

		person, ok := v.Get(0)
		require.True(t, ok)
		assert.Equal(t, "Alice", person.Name)
		assert.Equal(t, 30, person.Age)
	})
}

func TestVector_EdgeCases(t *testing.T) {
	t.Run("граница tail (32 элемента)", func(t *testing.T) {
		v := NewVector[int]()

		for i := 0; i < 32; i++ {
			v = v.Append(i)
		}
		assert.Equal(t, 32, v.Len())

		v = v.Append(32)
		assert.Equal(t, 33, v.Len(), "после переполнения tail длина должна быть 33")

		val, ok := v.Get(32)
		require.True(t, ok)
		assert.Equal(t, 32, val)
	})

	t.Run("несколько сбросов tail в дерево", func(t *testing.T) {
		v := NewVector[int]()

		for i := 0; i < 100; i++ {
			v = v.Append(i)
		}

		for i := 0; i < 100; i++ {
			val, ok := v.Get(i)
			require.True(t, ok, "элемент %d должен существовать", i)
			assert.Equal(t, i, val, "элемент %d должен быть равен %d", i, i)
		}
	})

	t.Run("Set в дереве vs в tail", func(t *testing.T) {
		v := NewVector[int]()
		for i := 0; i < 50; i++ {
			v = v.Append(i)
		}

		v2 := v.Set(10, 999)
		val, _ := v2.Get(10)
		assert.Equal(t, 999, val, "изменение в дереве должно работать")

		v3 := v.Set(45, 888)
		val, _ = v3.Get(45)
		assert.Equal(t, 888, val, "изменение в tail должно работать")
	})
}
