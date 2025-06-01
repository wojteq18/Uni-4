package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	num := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10,
		11, 12, 13, 14, 15, 16, 17, 18, 19, 20,
		21, 22, 23, 24, 25, 26, 27, 28, 29, 30}

	for i := 0; i < len(num)-1; i++ {
		fmt.Print(num[i], ",")
	}
	fmt.Print(num[29])

	rand.Seed(time.Now().UnixNano())

	// Tasowanie
	rand.Shuffle(len(num), func(i, j int) {
		num[i], num[j] = num[j], num[i]
	})

	fmt.Print(" ")

	for i := 0; i < len(num)-1; i++ {
		fmt.Print(num[i], ",")
	}
	fmt.Print(num[29])
}
