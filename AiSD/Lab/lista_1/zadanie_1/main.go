package main

import (
	"fmt"
	"zadanie_1/queues"
)

func main() {
	FifoQueue := queues.NewFifoQueue[int]()
	LifoQueue := queues.NewLifoQueue[int]()

	for i := 0; i < 50; i++ {
		FifoQueue.Push(i)
		LifoQueue.Push(i)
	}

	for i := 0; i < 40; i++ {
		a := FifoQueue.Remove()
		fmt.Println("FifoQueue: ", a)
		b := LifoQueue.Remove()
		fmt.Println("LifoQueue: ", b)
	}
}
