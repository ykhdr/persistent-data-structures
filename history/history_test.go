package history

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/ykhdr/persistent-data-structures/array"
)

func TestHistory_Basic(t *testing.T) {
	v := array.NewVector[int]()
	h := NewHistory(v)

	h.Commit(h.Current().Append(1))
	h.Commit(h.Current().Append(2))
	h.Commit(h.Current().Append(3))

	assert.Equal(t, 3, h.Current().Len(), "после трёх коммитов должно быть 3 элемента")
}

func TestHistory_Undo(t *testing.T) {
	t.Run("один откат", func(t *testing.T) {
		v := array.NewVector[int]()
		h := NewHistory(v)
		h.Commit(h.Current().Append(1))
		h.Commit(h.Current().Append(2))
		h.Commit(h.Current().Append(3))

		prev, ok := h.Undo()

		require.True(t, ok, "Undo должен вернуть ok")
		assert.Equal(t, 2, prev.Len(), "после отката должно быть 2 элемента")
	})

	t.Run("несколько откатов", func(t *testing.T) {
		v := array.NewVector[int]()
		h := NewHistory(v)
		h.Commit(h.Current().Append(1))
		h.Commit(h.Current().Append(2))
		h.Commit(h.Current().Append(3))

		h.Undo()
		prev, ok := h.Undo()

		require.True(t, ok, "второй Undo должен вернуть ok")
		assert.Equal(t, 1, prev.Len(), "после двух откатов должен остаться 1 элемент")
	})

	t.Run("откат на начальное состояние", func(t *testing.T) {
		v := array.NewVector[int]()
		h := NewHistory(v)
		h.Commit(h.Current().Append(1))

		prev, ok := h.Undo()

		require.True(t, ok)
		assert.Equal(t, 0, prev.Len(), "после отката должно быть начальное пустое состояние")
	})
}

func TestHistory_Redo(t *testing.T) {
	t.Run("один повтор", func(t *testing.T) {
		v := array.NewVector[int]()
		h := NewHistory(v)
		h.Commit(h.Current().Append(1))
		h.Commit(h.Current().Append(2))
		h.Undo()

		next, ok := h.Redo()

		require.True(t, ok, "Redo должен вернуть ok")
		assert.Equal(t, 2, next.Len(), "после повтора должно быть 2 элемента")
	})

	t.Run("несколько повторов", func(t *testing.T) {
		v := array.NewVector[int]()
		h := NewHistory(v)
		h.Commit(h.Current().Append(1))
		h.Commit(h.Current().Append(2))
		h.Commit(h.Current().Append(3))
		h.Undo()
		h.Undo()

		h.Redo()
		next, ok := h.Redo()

		require.True(t, ok, "второй Redo должен вернуть ok")
		assert.Equal(t, 3, next.Len(), "после двух повторов должно быть 3 элемента")
	})
}

func TestHistory_UndoLimit(t *testing.T) {
	t.Run("откат из начального состояния", func(t *testing.T) {
		v := array.NewVector[int]()
		h := NewHistory(v)

		_, ok := h.Undo()

		assert.False(t, ok, "Undo из начального состояния должен вернуть false")
	})

	t.Run("откат после исчерпания истории", func(t *testing.T) {
		v := array.NewVector[int]()
		h := NewHistory(v)
		h.Commit(h.Current().Append(1))
		h.Undo()

		_, ok := h.Undo()

		assert.False(t, ok, "Undo после полного отката должен вернуть false")
	})
}

func TestHistory_RedoLimit(t *testing.T) {
	t.Run("повтор на последней версии", func(t *testing.T) {
		v := array.NewVector[int]()
		h := NewHistory(v)
		h.Commit(h.Current().Append(1))

		_, ok := h.Redo()

		assert.False(t, ok, "Redo на последней версии должен вернуть false")
	})

	t.Run("повтор из начального состояния", func(t *testing.T) {
		v := array.NewVector[int]()
		h := NewHistory(v)

		_, ok := h.Redo()

		assert.False(t, ok, "Redo из начального состояния должен вернуть false")
	})
}

func TestHistory_BranchOverwrite(t *testing.T) {
	t.Run("новый коммит после отката стирает будущую историю", func(t *testing.T) {
		v := array.NewVector[int]()
		h := NewHistory(v)

		h.Commit(h.Current().Append(1))
		h.Commit(h.Current().Append(2))
		h.Commit(h.Current().Append(3))

		h.Undo()
		h.Undo()

		h.Commit(h.Current().Append(100))

		assert.Equal(t, 3, h.VersionCount(), "должно быть 3 версии после перезаписи")

		val, ok := h.Current().Get(1)
		require.True(t, ok)
		assert.Equal(t, 100, val, "второй элемент должен быть 100")
	})

	t.Run("Redo невозможен после нового коммита", func(t *testing.T) {
		v := array.NewVector[int]()
		h := NewHistory(v)

		h.Commit(h.Current().Append(1))
		h.Commit(h.Current().Append(2))
		h.Undo()
		h.Commit(h.Current().Append(999))

		_, ok := h.Redo()

		assert.False(t, ok, "Redo после нового коммита должен вернуть false")
	})
}

func TestHistory_CanUndoRedo(t *testing.T) {
	t.Run("начальное состояние", func(t *testing.T) {
		v := array.NewVector[int]()
		h := NewHistory(v)

		assert.False(t, h.CanUndo(), "CanUndo должен быть false в начальном состоянии")
		assert.False(t, h.CanRedo(), "CanRedo должен быть false в начальном состоянии")
	})

	t.Run("после коммита", func(t *testing.T) {
		v := array.NewVector[int]()
		h := NewHistory(v)
		h.Commit(h.Current().Append(1))

		assert.True(t, h.CanUndo(), "CanUndo должен быть true после коммита")
		assert.False(t, h.CanRedo(), "CanRedo должен быть false на последней версии")
	})

	t.Run("после отката", func(t *testing.T) {
		v := array.NewVector[int]()
		h := NewHistory(v)
		h.Commit(h.Current().Append(1))
		h.Undo()

		assert.False(t, h.CanUndo(), "CanUndo должен быть false в начальном состоянии")
		assert.True(t, h.CanRedo(), "CanRedo должен быть true после отката")
	})

	t.Run("в середине истории", func(t *testing.T) {
		v := array.NewVector[int]()
		h := NewHistory(v)
		h.Commit(h.Current().Append(1))
		h.Commit(h.Current().Append(2))
		h.Commit(h.Current().Append(3))
		h.Undo()

		assert.True(t, h.CanUndo(), "CanUndo должен быть true в середине истории")
		assert.True(t, h.CanRedo(), "CanRedo должен быть true в середине истории")
	})
}

func TestHistory_VersionCount(t *testing.T) {
	t.Run("начальное состояние", func(t *testing.T) {
		v := array.NewVector[int]()
		h := NewHistory(v)

		assert.Equal(t, 1, h.VersionCount(), "начальная версия должна считаться")
	})

	t.Run("после нескольких коммитов", func(t *testing.T) {
		v := array.NewVector[int]()
		h := NewHistory(v)
		h.Commit(h.Current().Append(1))
		h.Commit(h.Current().Append(2))
		h.Commit(h.Current().Append(3))

		assert.Equal(t, 4, h.VersionCount(), "должно быть 4 версии (начальная + 3 коммита)")
	})

	t.Run("откат не меняет количество версий", func(t *testing.T) {
		v := array.NewVector[int]()
		h := NewHistory(v)
		h.Commit(h.Current().Append(1))
		h.Commit(h.Current().Append(2))
		h.Undo()

		assert.Equal(t, 3, h.VersionCount(), "откат не должен удалять версии")
	})
}

func TestHistory_Persistence(t *testing.T) {
	t.Run("откат не меняет текущее состояние данных", func(t *testing.T) {
		v := array.NewVector[int]()
		h := NewHistory(v)

		h.Commit(h.Current().Append(1))
		h.Commit(h.Current().Append(2))

		current := h.Current()
		h.Undo()

		assert.Equal(t, 2, current.Len(), "сохранённая ссылка не должна измениться")
		assert.Equal(t, 1, h.Current().Len(), "текущее состояние истории должно измениться")
	})
}
