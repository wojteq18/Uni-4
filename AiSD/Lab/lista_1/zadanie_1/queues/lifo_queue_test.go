package queues

import "testing"

func TestPushLifo(t *testing.T) {
	fifoQueue := NewFifoQueue[int]()
	fifoQueue.Push(1)
	fifoQueue.Push(2)
	fifoQueue.Push(5)

	if fifoQueue.size != 3 {
		t.Errorf("Expected size: %d, got: %d", 3, fifoQueue.size)
	}

	if fifoQueue.front.value != 1 {
		t.Errorf("Expected value: %d, got: %d", 1, fifoQueue.front.value)
	}
}

func TestRemoveLifo(t *testing.T) {
	fifoQueue := NewFifoQueue[int]()
	fifoQueue.Push(1)
	fifoQueue.Push(2)
	fifoQueue.Push(5)

	removed1 := fifoQueue.Remove()
	removed2 := fifoQueue.Remove()

	if removed1 != 1 {
		t.Errorf("Expected removed value: %d, got: %d", 1, removed1)
	}
	if removed2 != 2 {
		t.Errorf("Expected removed value: %d, got: %d", 2, removed2)
	}
}
