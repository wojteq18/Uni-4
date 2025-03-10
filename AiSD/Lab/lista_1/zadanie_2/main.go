package main

import (
	queue "zadanie_2/queue"
)

func main() {
	list1 := queue.UndirectionalList{}

	for i := 11; i < 22; i++ { //adding 10 elements to the list
		queue.Insert(&list1, i)
	}
	list2 := queue.UndirectionalList{}

	for i := 31; i < 42; i++ { //adding 10 elements to the list
		queue.Insert(&list2, i)
	}

	list3 := queue.Merge(list1, list2)
	for i := 0; i < 20; i++ {
		println(list3.Remove())
	}
}
