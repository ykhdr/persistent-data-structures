package hashmap

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/ykhdr/persistent-data-structures/array"
)

func TestHashMap_SetAndGet(t *testing.T) {
	t.Run("добавление и получение элементов", func(t *testing.T) {
		m := NewHashMap[string, int]()

		m = m.Set("one", 1)
		m = m.Set("two", 2)
		m = m.Set("three", 3)

		assert.Equal(t, 3, m.Len(), "после добавления 3 элементов длина должна быть 3")

		val, ok := m.Get("one")
		require.True(t, ok, "ключ 'one' должен существовать")
		assert.Equal(t, 1, val)

		val, ok = m.Get("two")
		require.True(t, ok, "ключ 'two' должен существовать")
		assert.Equal(t, 2, val)

		val, ok = m.Get("three")
		require.True(t, ok, "ключ 'three' должен существовать")
		assert.Equal(t, 3, val)
	})
}

func TestHashMap_GetMissing(t *testing.T) {
	t.Run("получение несуществующего ключа", func(t *testing.T) {
		m := NewHashMap[string, int]().Set("one", 1)

		_, ok := m.Get("two")

		assert.False(t, ok, "Get несуществующего ключа должен вернуть false")
	})
}

func TestHashMap_Overwrite(t *testing.T) {
	t.Run("перезапись существующего ключа", func(t *testing.T) {
		m1 := NewHashMap[string, int]().Set("key", 1)
		m2 := m1.Set("key", 100)

		val1, _ := m1.Get("key")
		val2, _ := m2.Get("key")

		assert.Equal(t, 1, val1, "оригинал должен содержать старое значение")
		assert.Equal(t, 100, val2, "новая версия должна содержать новое значение")
		assert.Equal(t, 1, m1.Len(), "длина не должна измениться после перезаписи")
		assert.Equal(t, 1, m2.Len(), "длина не должна измениться после перезаписи")
	})
}

func TestHashMap_Persistence(t *testing.T) {
	t.Run("Set создаёт новую версию", func(t *testing.T) {
		m1 := NewHashMap[string, int]()
		m2 := m1.Set("a", 1)
		m3 := m2.Set("b", 2)
		m4 := m3.Set("c", 3)

		assert.Equal(t, 0, m1.Len(), "m1 не должен измениться")
		assert.Equal(t, 1, m2.Len(), "m2 не должен измениться")
		assert.Equal(t, 2, m3.Len(), "m3 не должен измениться")
		assert.Equal(t, 3, m4.Len(), "m4 должен содержать 3 элемента")
	})

	t.Run("старые версии не содержат новых ключей", func(t *testing.T) {
		m1 := NewHashMap[string, int]()
		m2 := m1.Set("a", 1)
		m3 := m2.Set("b", 2)

		assert.False(t, m2.Contains("b"), "m2 не должен содержать ключ 'b'")
		assert.True(t, m3.Contains("b"), "m3 должен содержать ключ 'b'")
	})
}

func TestHashMap_Delete(t *testing.T) {
	t.Run("удаление существующего ключа", func(t *testing.T) {
		m := NewHashMap[string, int]().Set("a", 1).Set("b", 2).Set("c", 3)

		m2 := m.Delete("b")

		assert.Equal(t, 3, m.Len(), "оригинал не должен измениться")
		assert.Equal(t, 2, m2.Len(), "новая версия должна содержать 2 элемента")
		assert.True(t, m.Contains("b"), "оригинал должен содержать 'b'")
		assert.False(t, m2.Contains("b"), "новая версия не должна содержать 'b'")
	})

	t.Run("остальные ключи сохраняются после удаления", func(t *testing.T) {
		m := NewHashMap[string, int]().Set("a", 1).Set("b", 2).Set("c", 3)
		m2 := m.Delete("b")

		val, ok := m2.Get("a")
		require.True(t, ok, "ключ 'a' должен остаться")
		assert.Equal(t, 1, val)

		val, ok = m2.Get("c")
		require.True(t, ok, "ключ 'c' должен остаться")
		assert.Equal(t, 3, val)
	})

	t.Run("удаление несуществующего ключа", func(t *testing.T) {
		m := NewHashMap[string, int]().Set("a", 1)

		m2 := m.Delete("nonexistent")

		assert.Equal(t, 1, m2.Len(), "удаление несуществующего ключа не должно менять длину")
	})
}

func TestHashMap_IntKeys(t *testing.T) {
	t.Run("целочисленные ключи", func(t *testing.T) {
		m := NewHashMap[int, string]()

		for i := 0; i < 100; i++ {
			m = m.Set(i, string(rune('a'+i%26)))
		}

		assert.Equal(t, 100, m.Len())

		for i := 0; i < 100; i++ {
			val, ok := m.Get(i)
			require.True(t, ok, "ключ %d должен существовать", i)
			expected := string(rune('a' + i%26))
			assert.Equal(t, expected, val, "значение для ключа %d должно быть '%s'", i, expected)
		}
	})
}

func TestHashMap_LargeDataset(t *testing.T) {
	t.Run("большой набор данных", func(t *testing.T) {
		m := NewHashMap[int, int]()

		for i := 0; i < 10000; i++ {
			m = m.Set(i, i*2)
		}

		assert.Equal(t, 10000, m.Len())

		for i := 0; i < 10000; i++ {
			val, ok := m.Get(i)
			require.True(t, ok, "ключ %d должен существовать", i)
			assert.Equal(t, i*2, val, "значение для ключа %d должно быть %d", i, i*2)
		}
	})
}

func TestHashMap_Iterator(t *testing.T) {
	t.Run("итерация по всем парам", func(t *testing.T) {
		m := NewHashMap[string, int]().Set("a", 1).Set("b", 2).Set("c", 3)

		count := 0
		sum := 0
		for _, v := range m.All() {
			count++
			sum += v
		}

		assert.Equal(t, 3, count, "должно быть 3 пары")
		assert.Equal(t, 6, sum, "сумма значений должна быть 6")
	})
}

func TestHashMap_Keys(t *testing.T) {
	t.Run("итерация по ключам", func(t *testing.T) {
		m := NewHashMap[string, int]().Set("a", 1).Set("b", 2)

		keys := make(map[string]bool)
		for k := range m.Keys() {
			keys[k] = true
		}

		assert.True(t, keys["a"], "ключ 'a' должен присутствовать")
		assert.True(t, keys["b"], "ключ 'b' должен присутствовать")
		assert.Equal(t, 2, len(keys), "должно быть 2 ключа")
	})
}

func TestHashMap_Values(t *testing.T) {
	t.Run("итерация по значениям", func(t *testing.T) {
		m := NewHashMap[string, int]().Set("a", 1).Set("b", 2)

		sum := 0
		count := 0
		for v := range m.Values() {
			sum += v
			count++
		}

		assert.Equal(t, 2, count, "должно быть 2 значения")
		assert.Equal(t, 3, sum, "сумма значений должна быть 3")
	})
}

func TestHashMap_NestedStructures(t *testing.T) {
	t.Run("хранение векторов в качестве значений", func(t *testing.T) {
		inner := array.NewVector[int]().Append(1).Append(2).Append(3)
		m := NewHashMap[string, *array.Vector[int]]().Set("numbers", inner)

		retrieved, ok := m.Get("numbers")
		require.True(t, ok, "ключ 'numbers' должен существовать")

		val, _ := retrieved.Get(0)
		assert.Equal(t, 1, val, "первый элемент вложенного вектора должен быть 1")
		assert.Equal(t, 3, retrieved.Len(), "вложенный вектор должен содержать 3 элемента")
	})
}

func TestHashMap_Contains(t *testing.T) {
	t.Run("проверка наличия ключа", func(t *testing.T) {
		m := NewHashMap[string, int]().Set("exists", 1)

		assert.True(t, m.Contains("exists"), "Contains должен вернуть true для существующего ключа")
		assert.False(t, m.Contains("missing"), "Contains должен вернуть false для отсутствующего ключа")
	})
}
