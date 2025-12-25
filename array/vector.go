package array

import "iter"

const (
	shiftStep = 5  // бит на уровень
	nodeWidth = 32 // 2^5 = 32 - количество значений в одной ноде дерева
	indexMask = 31 // (0b11111) - маска для извлечения 5 бит
)

type vectorNode[T any] struct {
	children [nodeWidth]*vectorNode[T] // значения для внутренних узлов
	values   [nodeWidth]T              // значения в листовом узле
}

func (n *vectorNode[T]) cloneInternal() *vectorNode[T] {
	newNode := &vectorNode[T]{}
	newNode.children = n.children
	return newNode
}

func (n *vectorNode[T]) cloneLeaf() *vectorNode[T] {
	newNode := &vectorNode[T]{}
	newNode.values = n.values
	return newNode
}

type Vector[T any] struct {
	root  *vectorNode[T] // корень дерева
	tail  []T            // буфер последних элементов (оптимизация)
	len   int            // количество элементов в структуре
	shift uint           // глубина дерева x 5 - нужна для побитового сдвига
}

func NewVector[T any]() *Vector[T] {
	return &Vector[T]{
		root:  nil,
		tail:  make([]T, 0, nodeWidth),
		len:   0,
		shift: shiftStep,
	}
}

func (v *Vector[T]) sealed() {}

func (v *Vector[T]) Len() int {
	return v.len
}

func (v *Vector[T]) tailOffset() int {
	if v.len < nodeWidth {
		return 0
	}
	return ((v.len - 1) >> shiftStep) << shiftStep
}

func (v *Vector[T]) getLeaf(index int) *vectorNode[T] {
	node := v.root
	for level := v.shift; level > shiftStep; level -= shiftStep {
		node = node.children[(index>>level)&indexMask]
	}
	return node.children[(index>>shiftStep)&indexMask]
}

func (v *Vector[T]) Get(index int) (T, bool) {
	var zero T
	if index < 0 || index >= v.len {
		return zero, false
	}

	if index >= v.tailOffset() {
		return v.tail[index-v.tailOffset()], true
	}

	leaf := v.getLeaf(index)
	return leaf.values[index&indexMask], true
}

func (v *Vector[T]) Set(index int, value T) *Vector[T] {
	if index < 0 || index >= v.len {
		return v
	}

	if index >= v.tailOffset() {
		newTail := make([]T, len(v.tail))
		copy(newTail, v.tail)
		newTail[index-v.tailOffset()] = value
		return &Vector[T]{
			root:  v.root,
			tail:  newTail,
			len:   v.len,
			shift: v.shift,
		}
	}

	return &Vector[T]{
		root:  v.setInNode(v.root, v.shift, index, value),
		tail:  v.tail,
		len:   v.len,
		shift: v.shift,
	}
}

func (v *Vector[T]) setInNode(node *vectorNode[T], level uint, index int, value T) *vectorNode[T] {
	if level == shiftStep {
		newNode := node.cloneInternal()
		childIndex := (index >> shiftStep) & indexMask
		leaf := node.children[childIndex].cloneLeaf()
		leaf.values[index&indexMask] = value
		newNode.children[childIndex] = leaf
		return newNode
	}

	newNode := node.cloneInternal()
	childIndex := (index >> level) & indexMask
	newNode.children[childIndex] = v.setInNode(node.children[childIndex], level-shiftStep, index, value)
	return newNode
}

func (v *Vector[T]) Append(value T) *Vector[T] {
	if len(v.tail) < nodeWidth {
		newTail := make([]T, len(v.tail)+1)
		copy(newTail, v.tail)
		newTail[len(v.tail)] = value
		return &Vector[T]{
			root:  v.root,
			tail:  newTail,
			len:   v.len + 1,
			shift: v.shift,
		}
	}

	tailNode := &vectorNode[T]{}
	copy(tailNode.values[:], v.tail)

	var newRoot *vectorNode[T]
	newShift := v.shift

	if v.root == nil {
		newRoot = &vectorNode[T]{}
		newRoot.children[0] = tailNode
	} else if (v.len >> shiftStep) > (1 << v.shift) {
		newRoot = &vectorNode[T]{}
		newRoot.children[0] = v.root
		newRoot.children[1] = v.newPath(v.shift, tailNode)
		newShift += shiftStep
	} else {
		newRoot = v.pushTail(v.shift, v.root, tailNode)
	}

	return &Vector[T]{
		root:  newRoot,
		tail:  []T{value},
		len:   v.len + 1,
		shift: newShift,
	}
}

func (v *Vector[T]) newPath(level uint, leaf *vectorNode[T]) *vectorNode[T] {
	if level == shiftStep {
		node := &vectorNode[T]{}
		node.children[0] = leaf
		return node
	}
	node := &vectorNode[T]{}
	node.children[0] = v.newPath(level-shiftStep, leaf)
	return node
}

func (v *Vector[T]) pushTail(level uint, parent *vectorNode[T], tailNode *vectorNode[T]) *vectorNode[T] {
	subIndex := ((v.len - 1) >> level) & indexMask
	newNode := parent.cloneInternal()

	if level == shiftStep {
		newNode.children[subIndex] = tailNode
	} else {
		child := parent.children[subIndex]
		if child != nil {
			newNode.children[subIndex] = v.pushTail(level-shiftStep, child, tailNode)
		} else {
			newNode.children[subIndex] = v.newPath(level-shiftStep, tailNode)
		}
	}

	return newNode
}

func (v *Vector[T]) Pop() (*Vector[T], T, bool) {
	var zero T
	if v.len == 0 {
		return v, zero, false
	}

	if v.len == 1 {
		value := v.tail[0]
		return NewVector[T](), value, true
	}

	if len(v.tail) > 1 {
		newTail := make([]T, len(v.tail)-1)
		copy(newTail, v.tail[:len(v.tail)-1])
		value := v.tail[len(v.tail)-1]
		return &Vector[T]{
			root:  v.root,
			tail:  newTail,
			len:   v.len - 1,
			shift: v.shift,
		}, value, true
	}

	value := v.tail[0]
	newTail := v.leafValuesToSlice(v.len - 2)

	var newRoot *vectorNode[T]
	newShift := v.shift

	newRoot = v.popTail(v.shift, v.root)
	if newRoot != nil && newRoot.children[1] == nil && v.shift > shiftStep {
		newRoot = newRoot.children[0]
		newShift -= shiftStep
	}

	return &Vector[T]{
		root:  newRoot,
		tail:  newTail,
		len:   v.len - 1,
		shift: newShift,
	}, value, true
}

func (v *Vector[T]) leafValuesToSlice(index int) []T {
	leaf := v.getLeaf(index)
	result := make([]T, nodeWidth)
	copy(result, leaf.values[:])
	return result
}

func (v *Vector[T]) popTail(level uint, node *vectorNode[T]) *vectorNode[T] {
	subIndex := ((v.len - 2) >> level) & indexMask

	if level > shiftStep {
		newChild := v.popTail(level-shiftStep, node.children[subIndex])
		if newChild == nil && subIndex == 0 {
			return nil
		}
		newNode := node.cloneInternal()
		newNode.children[subIndex] = newChild
		return newNode
	}

	if subIndex == 0 {
		return nil
	}

	newNode := node.cloneInternal()
	newNode.children[subIndex] = nil
	return newNode
}

func (v *Vector[T]) All() iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		for i := 0; i < v.len; i++ {
			value, _ := v.Get(i)
			if !yield(i, value) {
				return
			}
		}
	}
}

func (v *Vector[T]) Values() iter.Seq[T] {
	return func(yield func(T) bool) {
		for i := 0; i < v.len; i++ {
			value, _ := v.Get(i)
			if !yield(value) {
				return
			}
		}
	}
}
