package main

import (
	"fmt"
	"math/rand/v2"
	"zadanie_2/queue"
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
		fmt.Println(list3.Remove())
	}
	fmt.Println(rand.IntN(100), ",")
	var T [10000]int
	for i := 0; i < 10000; i++ {
		T[i] = rand.IntN(100000)
	}
	fmt.Println(T[100])
	L := queue.UndirectionalList{}
	for i := 0; i < 10000; i++ {
		queue.Insert(&L, T[i])
	}
	/*amount := 0
	for i := 0; i <= 1000; i++ {
		randNumber := rand.IntN(100000)
		for j := 0; j < 10000; j++ {
			amount++
			if queue.Cointains(L, T[randNumber]) {
				break
			}
		}

	}
	avg := amount / 1000
	fmt.Println(avg)*/
}
