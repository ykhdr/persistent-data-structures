package history

type History[T any] struct {
	versions []T
	current  int
}

func NewHistory[T any](initial T) *History[T] {
	return &History[T]{
		versions: []T{initial},
		current:  0,
	}
}

func (h *History[T]) Current() T {
	return h.versions[h.current]
}

func (h *History[T]) Commit(newVersion T) {
	h.versions = h.versions[:h.current+1]
	h.versions = append(h.versions, newVersion)
	h.current++
}

func (h *History[T]) Undo() (T, bool) {
	if h.current > 0 {
		h.current--
		return h.versions[h.current], true
	}
	var zero T
	return zero, false
}

func (h *History[T]) Redo() (T, bool) {
	if h.current < len(h.versions)-1 {
		h.current++
		return h.versions[h.current], true
	}
	var zero T
	return zero, false
}

func (h *History[T]) CanUndo() bool {
	return h.current > 0
}

func (h *History[T]) CanRedo() bool {
	return h.current < len(h.versions)-1
}

func (h *History[T]) VersionCount() int {
	return len(h.versions)
}

func (h *History[T]) CurrentIndex() int {
	return h.current
}
