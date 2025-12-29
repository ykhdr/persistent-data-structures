package hashmap

type NaiveHashMap[K comparable, V any] struct {
	data map[K]V
}

func NewNaiveHashMap[K comparable, V any]() *NaiveHashMap[K, V] {
	return &NaiveHashMap[K, V]{data: make(map[K]V)}
}

func (m *NaiveHashMap[K, V]) Len() int { return len(m.data) }

func (m *NaiveHashMap[K, V]) Get(key K) (V, bool) {
	v, ok := m.data[key]
	return v, ok
}

func (m *NaiveHashMap[K, V]) Set(key K, value V) *NaiveHashMap[K, V] {
	newData := make(map[K]V, len(m.data)+1)
	for k, v := range m.data {
		newData[k] = v
	}
	newData[key] = value
	return &NaiveHashMap[K, V]{data: newData}
}

func (m *NaiveHashMap[K, V]) Delete(key K) *NaiveHashMap[K, V] {
	if _, ok := m.data[key]; !ok {
		return m
	}
	newData := make(map[K]V, len(m.data))
	for k, v := range m.data {
		if k == key {
			continue
		}
		newData[k] = v
	}
	return &NaiveHashMap[K, V]{data: newData}
}

func (m *NaiveHashMap[K, V]) Contains(key K) bool {
	_, ok := m.data[key]
	return ok
}
