package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

func sortNines(arr []int, s *int, c *int) {
	for i := 1; i < len(arr); i++ {
		key := arr[i]
		j := i - 1
		for j >= 0 {
			(*c)++
			if arr[j] > key {
				arr[j+1] = arr[j]
				j--
			} else {
				break
			}
		}
		arr[j+1] = key
	}
}

func findMedianInNine(arr []int, s *int, c *int) int {
	sortNines(arr, s, c)
	length := len(arr)

	if length == 0 {
		panic("Array is empty")
	}

	return arr[length/2]
}

func medianOfMedians(arr []int, s *int, c *int) int {

	length := len(arr)

	if length <= 9 {
		return findMedianInNine(arr, s, c)
	}

	var medians []int
	for i := 0; i < length; i += 9 {
		end := i + 9
		if end > length {
			end = length
		}
		slice := arr[i:end]
		medians = append(medians, findMedianInNine(slice, s, c))
	}

	return selectFind(medians, len(medians)/2+1, s, c)
}

func selectFind(arr []int, k int, s *int, c *int) int {
	length := len(arr)
	pivotValue := medianOfMedians(arr, s, c)

	idx := 0
	for i := 0; i < length; i++ {
		(*c)++ // Porównanie wartości
		if arr[i] == pivotValue {
			idx = i
			break
		}
	}

	if idx != length-1 {
		arr[idx], arr[length-1] = arr[length-1], arr[idx]
		(*s)++ // swap
	}

	i := 0
	for j := 0; j < length-1; j++ {
		(*c)++ // Porównanie wartości
		if arr[j] < arr[length-1] {
			if i != j {
				arr[i], arr[j] = arr[j], arr[i]
				(*s)++ // swap
			}
			i++
		}
	}

	if i != length-1 {
		arr[i], arr[length-1] = arr[length-1], arr[i]
		(*s)++ // swap
	}

	if i < k-1 {
		right := arr[i+1:]
		return selectFind(right, k-i-1, s, c)
	}

	if i > k-1 {
		left := arr[:i]
		return selectFind(left, k, s, c)
	}

	return arr[i]
}

func main() {
	s := 0
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

	startTime := time.Now()
	a := selectFind(numbers, n, &s, &comp)
	c := time.Since(startTime).Microseconds()

	fmt.Println("s =", s)
	fmt.Println("c =", c)
	fmt.Println("znaleziony element:", a)
}
