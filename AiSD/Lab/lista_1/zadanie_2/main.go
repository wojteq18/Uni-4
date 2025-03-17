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

	list3 := queue.Merge(&list1, &list2)
	for i := 0; i < 20; i++ {
		fmt.Println("l3: ", list3.Remove())
	}

	var T [10000]int //array with 10000 random numbers
	for i := 0; i < 10000; i++ {
		T[i] = rand.IntN(100000)
	}
	//fmt.Println(T[100])
	L := queue.UndirectionalList{} //list with 10000 random numbers
	for i := 0; i < 10000; i++ {
		queue.Insert(&L, T[i])
	}
	all1 := 0
	for i := 0; i < 10000; i++ { //case 1 - random number from T
		randNumber := rand.IntN(10000)
		_, comps := queue.Contains(L, T[randNumber])
		all1 += comps
	}
	fmt.Println("Srednia 1 : ", all1/10000)

	all2 := 0
	for i := 0; i < 10000; i++ { //case 2 - random number from 0 to 100000
		randNumber := rand.IntN(100000)
		_, comps := queue.Contains(L, randNumber)
		all2 += comps
	}
	fmt.Println("Srednia 2 : ", all2/10000)
}
