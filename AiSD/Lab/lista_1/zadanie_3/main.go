package main

import (
	"fmt"
	"math/rand/v2"
	"zadanie_3/queue"
)

func main() {
	list1 := queue.CircularDll{}
	list2 := queue.CircularDll{}

	for i := 11; i < 21; i++ {
		queue.Insert(&list1, i)
	}

	for i := 31; i < 41; i++ {
		queue.Insert(&list2, i)
	}

	list3 := queue.Merge(&list1, &list2)
	fmt.Println("hejla")
	for i := 0; i < 20; i++ {
		fmt.Println("element z l3: ", list3.Remove())
	}

	var T [10000]int
	for i := 0; i < 10000; i++ {
		T[i] = rand.IntN(100000)
	}

	L := queue.CircularDll{} //list with 10000 random numbers
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
