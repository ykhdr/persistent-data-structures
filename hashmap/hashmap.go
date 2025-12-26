package hashmap

import (
	"hash/maphash"
	"iter"
	"math/bits"
)

const (
	hmapShift = 5
	hmapMask  = 31
)

type entry[K comparable, V any] struct {
	key   K
	value V
}

type collision[K comparable, V any] struct {
	entries []entry[K, V]
}

type hmapNode[K comparable, V any] struct {
	bitmap   uint32
	children []any
}

func (n *hmapNode[K, V]) index(bit uint32) int {
	return bits.OnesCount32(n.bitmap & (bit - 1))
}

func (n *hmapNode[K, V]) clone() *hmapNode[K, V] {
	newChildren := make([]any, len(n.children))
	copy(newChildren, n.children)
	return &hmapNode[K, V]{
		bitmap:   n.bitmap,
		children: newChildren,
	}
}

type HashMap[K comparable, V any] struct {
	root *hmapNode[K, V]
	len  int
	seed maphash.Seed
}

func NewHashMap[K comparable, V any]() *HashMap[K, V] {
	return &HashMap[K, V]{
		root: &hmapNode[K, V]{},
		len:  0,
		seed: maphash.MakeSeed(),
	}
}

func (m *HashMap[K, V]) Len() int {
	return m.len
}

func (m *HashMap[K, V]) hash(key K) uint32 {
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
	return uint32(h.Sum64())
}

func (m *HashMap[K, V]) Get(key K) (V, bool) {
	var zero V
	if m.root == nil {
		return zero, false
	}

	hash := m.hash(key)
	return m.getNode(m.root, key, hash, 0)
}

func (m *HashMap[K, V]) getNode(node *hmapNode[K, V], key K, hash uint32, shift uint) (V, bool) {
	var zero V
	bit := uint32(1) << ((hash >> shift) & hmapMask)

	if node.bitmap&bit == 0 {
		return zero, false
	}

	idx := node.index(bit)
	child := node.children[idx]

	switch c := child.(type) {
	case *hmapNode[K, V]:
		return m.getNode(c, key, hash, shift+hmapShift)
	case *entry[K, V]:
		if c.key == key {
			return c.value, true
		}
		return zero, false
	case *collision[K, V]:
		for _, e := range c.entries {
			if e.key == key {
				return e.value, true
			}
		}
		return zero, false
	}

	return zero, false
}

func (m *HashMap[K, V]) Set(key K, value V) *HashMap[K, V] {
	hash := m.hash(key)
	newRoot, added := m.setNode(m.root, key, value, hash, 0)

	newLen := m.len
	if added {
		newLen++
	}

	return &HashMap[K, V]{
		root: newRoot,
		len:  newLen,
		seed: m.seed,
	}
}

func (m *HashMap[K, V]) setNode(node *hmapNode[K, V], key K, value V, hash uint32, shift uint) (*hmapNode[K, V], bool) {
	bit := uint32(1) << ((hash >> shift) & hmapMask)
	idx := node.index(bit)

	if node.bitmap&bit == 0 {
		newNode := node.clone()
		newNode.bitmap |= bit
		newChildren := make([]any, len(node.children)+1)
		copy(newChildren[:idx], node.children[:idx])
		newChildren[idx] = &entry[K, V]{key: key, value: value}
		copy(newChildren[idx+1:], node.children[idx:])
		newNode.children = newChildren
		return newNode, true
	}

	child := node.children[idx]
	newNode := node.clone()

	switch c := child.(type) {
	case *hmapNode[K, V]:
		newChild, added := m.setNode(c, key, value, hash, shift+hmapShift)
		newNode.children[idx] = newChild
		return newNode, added

	case *entry[K, V]:
		if c.key == key {
			newNode.children[idx] = &entry[K, V]{key: key, value: value}
			return newNode, false
		}

		existingHash := m.hash(c.key)

		if shift >= 30 {
			newNode.children[idx] = &collision[K, V]{
				entries: []entry[K, V]{
					{key: c.key, value: c.value},
					{key: key, value: value},
				},
			}
			return newNode, true
		}

		newChild := m.createTwoEntryNode(c.key, c.value, existingHash, key, value, hash, shift+hmapShift)
		newNode.children[idx] = newChild
		return newNode, true

	case *collision[K, V]:
		for i, e := range c.entries {
			if e.key == key {
				newEntries := make([]entry[K, V], len(c.entries))
				copy(newEntries, c.entries)
				newEntries[i] = entry[K, V]{key: key, value: value}
				newNode.children[idx] = &collision[K, V]{entries: newEntries}
				return newNode, false
			}
		}
		newEntries := make([]entry[K, V], len(c.entries)+1)
		copy(newEntries, c.entries)
		newEntries[len(c.entries)] = entry[K, V]{key: key, value: value}
		newNode.children[idx] = &collision[K, V]{entries: newEntries}
		return newNode, true
	}

	return newNode, false
}

func (m *HashMap[K, V]) createTwoEntryNode(key1 K, val1 V, hash1 uint32, key2 K, val2 V, hash2 uint32, shift uint) *hmapNode[K, V] {
	if shift >= 30 {
		return &hmapNode[K, V]{
			bitmap: 1,
			children: []any{
				&collision[K, V]{
					entries: []entry[K, V]{
						{key: key1, value: val1},
						{key: key2, value: val2},
					},
				},
			},
		}
	}

	bit1 := uint32(1) << ((hash1 >> shift) & hmapMask)
	bit2 := uint32(1) << ((hash2 >> shift) & hmapMask)

	if bit1 == bit2 {
		child := m.createTwoEntryNode(key1, val1, hash1, key2, val2, hash2, shift+hmapShift)
		return &hmapNode[K, V]{
			bitmap:   bit1,
			children: []any{child},
		}
	}

	if bit1 < bit2 {
		return &hmapNode[K, V]{
			bitmap: bit1 | bit2,
			children: []any{
				&entry[K, V]{key: key1, value: val1},
				&entry[K, V]{key: key2, value: val2},
			},
		}
	}

	return &hmapNode[K, V]{
		bitmap: bit1 | bit2,
		children: []any{
			&entry[K, V]{key: key2, value: val2},
			&entry[K, V]{key: key1, value: val1},
		},
	}
}

func (m *HashMap[K, V]) Delete(key K) *HashMap[K, V] {
	if m.root == nil {
		return m
	}

	hash := m.hash(key)
	newRoot, deleted := m.deleteNode(m.root, key, hash, 0)

	if !deleted {
		return m
	}

	return &HashMap[K, V]{
		root: newRoot,
		len:  m.len - 1,
		seed: m.seed,
	}
}

func (m *HashMap[K, V]) deleteNode(node *hmapNode[K, V], key K, hash uint32, shift uint) (*hmapNode[K, V], bool) {
	bit := uint32(1) << ((hash >> shift) & hmapMask)

	if node.bitmap&bit == 0 {
		return node, false
	}

	idx := node.index(bit)
	child := node.children[idx]

	switch c := child.(type) {
	case *hmapNode[K, V]:
		newChild, deleted := m.deleteNode(c, key, hash, shift+hmapShift)
		if !deleted {
			return node, false
		}

		newNode := node.clone()

		if newChild.bitmap == 0 {
			newNode.bitmap &^= bit
			newChildren := make([]any, len(node.children)-1)
			copy(newChildren[:idx], node.children[:idx])
			copy(newChildren[idx:], node.children[idx+1:])
			newNode.children = newChildren
		} else if len(newChild.children) == 1 {
			if e, ok := newChild.children[0].(*entry[K, V]); ok {
				newNode.children[idx] = e
			} else {
				newNode.children[idx] = newChild
			}
		} else {
			newNode.children[idx] = newChild
		}

		return newNode, true

	case *entry[K, V]:
		if c.key != key {
			return node, false
		}

		newNode := node.clone()
		newNode.bitmap &^= bit
		newChildren := make([]any, len(node.children)-1)
		copy(newChildren[:idx], node.children[:idx])
		copy(newChildren[idx:], node.children[idx+1:])
		newNode.children = newChildren

		return newNode, true

	case *collision[K, V]:
		foundIdx := -1
		for i, e := range c.entries {
			if e.key == key {
				foundIdx = i
				break
			}
		}

		if foundIdx == -1 {
			return node, false
		}

		newNode := node.clone()

		if len(c.entries) == 2 {
			remaining := c.entries[1-foundIdx]
			newNode.children[idx] = &entry[K, V]{key: remaining.key, value: remaining.value}
		} else {
			newEntries := make([]entry[K, V], len(c.entries)-1)
			copy(newEntries[:foundIdx], c.entries[:foundIdx])
			copy(newEntries[foundIdx:], c.entries[foundIdx+1:])
			newNode.children[idx] = &collision[K, V]{entries: newEntries}
		}

		return newNode, true
	}

	return node, false
}

func (m *HashMap[K, V]) Contains(key K) bool {
	_, ok := m.Get(key)
	return ok
}

func (m *HashMap[K, V]) All() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		if m.root == nil {
			return
		}
		m.iterNode(m.root, yield)
	}
}

func (m *HashMap[K, V]) iterNode(node *hmapNode[K, V], yield func(K, V) bool) bool {
	for _, child := range node.children {
		switch c := child.(type) {
		case *hmapNode[K, V]:
			if !m.iterNode(c, yield) {
				return false
			}
		case *entry[K, V]:
			if !yield(c.key, c.value) {
				return false
			}
		case *collision[K, V]:
			for _, e := range c.entries {
				if !yield(e.key, e.value) {
					return false
				}
			}
		}
	}
	return true
}

func (m *HashMap[K, V]) Keys() iter.Seq[K] {
	return func(yield func(K) bool) {
		for k, _ := range m.All() {
			if !yield(k) {
				return
			}
		}
	}
}

func (m *HashMap[K, V]) Values() iter.Seq[V] {
	return func(yield func(V) bool) {
		for _, v := range m.All() {
			if !yield(v) {
				return
			}
		}
	}
}
