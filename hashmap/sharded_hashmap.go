package hashmap

import (
	"hash/maphash"
	"iter"
)

const defaultBuckets = 256

type ShardedHashMap[K comparable, V any] struct {
	buckets []map[K]V
	len     int
	seed    maphash.Seed
}

func NewShardedHashMap[K comparable, V any]() *ShardedHashMap[K, V] {
	return &ShardedHashMap[K, V]{
		buckets: make([]map[K]V, defaultBuckets),
		len:     0,
		seed:    maphash.MakeSeed(),
	}
}

func (m *ShardedHashMap[K, V]) Len() int {
	return m.len
}

func (m *ShardedHashMap[K, V]) hash(key K) uint64 {
	var h maphash.Hash
	h.SetSeed(m.seed)

	switch k := any(key).(type) {
	case string:
		h.WriteString(k)
	case int:
		var buf [8]byte
		for i := 0; i < 8; i++ {
			buf[i] = byte(k >> (i * 8))
		}
		h.Write(buf[:])
	case int64:
		var buf [8]byte
		for i := 0; i < 8; i++ {
			buf[i] = byte(k >> (i * 8))
		}
		h.Write(buf[:])
	case int32:
		var buf [4]byte
		for i := 0; i < 4; i++ {
			buf[i] = byte(k >> (i * 8))
		}
		h.Write(buf[:])
	case uint:
		var buf [8]byte
		for i := 0; i < 8; i++ {
			buf[i] = byte(k >> (i * 8))
		}
		h.Write(buf[:])
	case uint64:
		var buf [8]byte
		for i := 0; i < 8; i++ {
			buf[i] = byte(k >> (i * 8))
		}
		h.Write(buf[:])
	case uint32:
		var buf [4]byte
		for i := 0; i < 4; i++ {
			buf[i] = byte(k >> (i * 8))
		}
		h.Write(buf[:])
	default:
		h.WriteString(any(key).(string))
	}

	return h.Sum64()
}

func (m *ShardedHashMap[K, V]) bucketIndex(key K) int {
	return int(m.hash(key) % uint64(len(m.buckets)))
}

func (m *ShardedHashMap[K, V]) Get(key K) (V, bool) {
	var zero V
	if len(m.buckets) == 0 {
		return zero, false
	}

	i := m.bucketIndex(key)
	b := m.buckets[i]
	if b == nil {
		return zero, false
	}

	v, ok := b[key]
	return v, ok
}

func (m *ShardedHashMap[K, V]) Contains(key K) bool {
	_, ok := m.Get(key)
	return ok
}

func (m *ShardedHashMap[K, V]) Set(key K, value V) *ShardedHashMap[K, V] {
	i := m.bucketIndex(key)
	oldB := m.buckets[i]

	_, existed := oldB[key]

	newBuckets := make([]map[K]V, len(m.buckets))
	copy(newBuckets, m.buckets)

	newB := make(map[K]V, len(oldB)+1)
	for k, v := range oldB {
		newB[k] = v
	}
	newB[key] = value
	newBuckets[i] = newB

	newLen := m.len
	if !existed {
		newLen++
	}

	return &ShardedHashMap[K, V]{
		buckets: newBuckets,
		len:     newLen,
		seed:    m.seed,
	}
}

func (m *ShardedHashMap[K, V]) Delete(key K) *ShardedHashMap[K, V] {
	i := m.bucketIndex(key)
	oldB := m.buckets[i]
	if oldB == nil {
		return m
	}

	if _, ok := oldB[key]; !ok {
		return m
	}

	newBuckets := make([]map[K]V, len(m.buckets))
	copy(newBuckets, m.buckets)

	if len(oldB) == 1 {
		// после удаления бакет станет пустым
		newBuckets[i] = nil
	} else {
		newB := make(map[K]V, len(oldB)-1)
		for k, v := range oldB {
			if k == key {
				continue
			}
			newB[k] = v
		}
		newBuckets[i] = newB
	}

	return &ShardedHashMap[K, V]{
		buckets: newBuckets,
		len:     m.len - 1,
		seed:    m.seed,
	}
}

func (m *ShardedHashMap[K, V]) All() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for _, b := range m.buckets {
			if b == nil {
				continue
			}
			for k, v := range b {
				if !yield(k, v) {
					return
				}
			}
		}
	}
}

func (m *ShardedHashMap[K, V]) Keys() iter.Seq[K] {
	return func(yield func(K) bool) {
		for k := range m.All() {
			if !yield(k) {
				return
			}
		}
	}
}

func (m *ShardedHashMap[K, V]) Values() iter.Seq[V] {
	return func(yield func(V) bool) {
		for _, v := range m.All() {
			if !yield(v) {
				return
			}
		}
	}
}
