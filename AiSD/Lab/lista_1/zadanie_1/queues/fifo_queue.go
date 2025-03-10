package queues

type NodeFifo[T any] struct { //Node represents a single element in queue
	value T
	next  *NodeFifo[T]
}

type FifoQueue[T any] struct { //represents fifo queue
	front *NodeFifo[T]
	back  *NodeFifo[T]
	size  int
}

func NewFifoQueue[T any]() *FifoQueue[T] { //creates new fifo queue
	return &FifoQueue[T]{
		front: nil,
		back:  nil,
		size:  0,
	}
}

func (q *FifoQueue[T]) Push(value T) { //pushes element to fifo queue
	newNode := &NodeFifo[T]{value: value, next: nil}
	if q.back != nil {
		q.back.next = newNode //if back is not nil, then set next of back to new node
	} else {
		q.front = newNode
	}
	q.back = newNode
	q.size++
}

func (q *FifoQueue[T]) Remove() T { //removes element from fifo queue
	if q.size == 0 {
		panic("Queue is empty")
	} else {
		value := q.front.value
		q.front = q.front.next
		if q.front == nil {
			q.back = nil
		}
		q.size--
		return value
	}
}
