package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

func binarySearch(arr []int, target int, c *int) int {
	length := len(arr)

	if length == 0 {
		return -1
	}

	if length == 1 {
		(*c)++
		if arr[0] == target {
			return 1
		} else {
			return -1
		}
	}

	pivot := arr[length/2]
	//fmt.Println("Aktualny stan: ", arr)

	(*c)++
	if pivot == target {
		return 1
	} else if target > pivot {
		(*c)++
		return binarySearch(arr[length/2:], target, c)
	} else {
		(*c)++
		return binarySearch(arr[:length/2], target, c)
	}
}

func main() {
	comp := 0

	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintln(os.Stderr, "I/O Error reading stdin")
		return
	}

	line := string(input)
	start := strings.Index(line, "[")
	end := strings.Index(line, "]")
	if start == -1 || end == -1 || end <= start {
		fmt.Fprintln(os.Stderr, "I/O Error parsing array")
		return
	}

	nStr := strings.TrimSpace(line[end+1:])
	n, err := strconv.Atoi(nStr)
	if err != nil {
		fmt.Fprintln(os.Stderr, "n is not correct")
		return
	}

	numbersStr := line[start+1 : end]
	parts := strings.Split(numbersStr, ",")

	var numbers []int
	for _, part := range parts {
		num, err := strconv.Atoi(strings.TrimSpace(part))
		if err != nil {
			fmt.Fprintln(os.Stderr, "parse error")
			return
		}
		numbers = append(numbers, num)
	}

	//fmt.Println("Pierwotna tablica: ", numbers)
	//fmt.Println("Szukana liczba: ", n)
	startTime := time.Now()
	binarySearch(numbers, n, &comp)
	c := time.Since(startTime)

	//fmt.Println("Rezultat: ", result)
	fmt.Println("c =", c)
}

//liczba porównań (O(1)) ~ 2
