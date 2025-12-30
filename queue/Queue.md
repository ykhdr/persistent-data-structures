# 3. Persistent Queue

**Описание**

Persistent Queue — неизменяемая очередь FIFO (first-in-first-out).
Реализована на основе **двух persistent-стеков** (`front` и `rear`) и операции разворота `rear.reverse()`,
которая выполняется только когда `front` пуст.

**Реализация**

- Основана на двух persistent-стеках:
    - `front` — для извлечения (dequeue / peek)
    - `rear` — для добавления (enqueue)
- При опустевшем `front` выполняется ленивый разворот `rear -> front` (reverse)
- Используется structural sharing: стек — persistent связный список, `push`/`pop` создают новые версии без копирования
  всей структуры

**Сложность операций**

| Операция        | Сложность                                     |
|-----------------|-----------------------------------------------|
| Enqueue         | $O(1)$                                        |
| Dequeue         | амортизированное $O(1)$, худший случай $O(n)$ |
| Peek            | амортизированное $O(1)$, худший случай $O(n)$ |
| Len / IsEmpty   | $O(1)$                                        |
| Iteration (All) | $O(n)$                                        |

## Архитектура

Очередь хранит элементы в двух стеках:

```go
type Queue[T any] struct {
  front *stack[T]
  rear  *stack[T]
  len   int
}
```

```go
type stackNode[T any] struct {
  value T
  next  *stackNode[T]
}
```

```go
type stack[T any] struct {
  head *stackNode[T]
  len int
}
```