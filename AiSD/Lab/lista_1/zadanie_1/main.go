package main

import (
	"fmt"
	"zadanie_1/queues"
)

func main() {
	fmt.Println("jjj")
	q := queues.NewFifoQueue[int]()
	q.Push(1)
	w := q.Remove()
	fmt.Println(w)
}
