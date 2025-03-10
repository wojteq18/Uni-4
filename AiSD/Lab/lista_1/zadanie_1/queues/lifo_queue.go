package queues

type NodeLifo[T any] struct {
	value T
	next  *NodeLifo[T]
}

type LifoQueue[T any] struct {
	top  *NodeLifo[T]
	size int
}

func NewLifoQueue[T any]() *LifoQueue[T] {
	return &LifoQueue[T]{
		top:  nil,
		size: 0,
	}
}

func (q *LifoQueue[T]) Push(value T) {
	newNode := &NodeLifo[T]{value: value, next: q.top}
	q.top = newNode
	q.size++
}

func (q *LifoQueue[T]) Remove() T {
	if q.size == 0 {
		panic("ligo queue is empty!")
	} else {
		value := q.top.value
		q.top = q.top.next
		q.size--
		return value
	}
}
