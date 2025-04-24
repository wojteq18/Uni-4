package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

func findMedianInFive(arr *[]int) int {
	sort.Ints(*arr)
	length := len(*arr)

	if length == 0 {
		panic("Array is empty")
	}

	return (*arr)[length/2]
}

func medianOfMedians(arr *[]int) int {

	length := len(*arr)

	if length <= 5 {
		return findMedianInFive(arr)
	}

	var medians []int
	for i := 0; i < length; i += 5 {
		end := i + 5
		if end > length {
			end = length
		}
		subArr := (*arr)[i:end]
		medians = append(medians, findMedianInFive(&subArr))
	}

	return medianOfMedians(&medians)
}

func selectFind(arr *[]int, k int) int {
	length := len(*arr)
	pivotValue := medianOfMedians(arr)

	idx := 0

	for i := 0; i < length; i++ {
		if (*arr)[i] == pivotValue {
			idx = i
			break
		}
	}

	if idx != length-1 {
		(*arr)[idx], (*arr)[length-1] = (*arr)[length-1], (*arr)[idx]
	}

	i := 0
	for j := 0; j < length-1; j++ {
		if (*arr)[j] < (*arr)[length-1] {
			if i != j {
				(*arr)[i], (*arr)[j] = (*arr)[j], (*arr)[i]
			}
			i++
		}
	}

	if i != length-1 {
		(*arr)[i], (*arr)[length-1] = (*arr)[length-1], (*arr)[i]
	}

	if i < k-1 {
		right := (*arr)[i+1:]
		return selectFind(&right, k-i-1)
	}

	if i > k-1 {
		left := (*arr)[:i]
		return selectFind(&left, k)
	}

	return (*arr)[i]
}

func main() {
	s := 0
	c := 0
	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		fmt.Fprintln(os.Stderr, "I/O Error")
		return
	}
	line := scanner.Text()

	start := strings.Index(line, "[")
	end := strings.Index(line, "]")
	if start == -1 || end == -1 || end <= start { // Jeśli nie znalazł nawiasu - zwraca -1
		fmt.Fprintln(os.Stderr, "I/O Error")
		return
	}

	nStr := strings.TrimSpace(line[end+1:])
	n, err := strconv.Atoi(nStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "n is not correct")
		return
	}

	// Parsowanie tablicy
	numbersStr := line[start+1 : end]
	parts := strings.Split(numbersStr, ",")

	var numbers []int
	for _, part := range parts {
		num, err := strconv.Atoi(strings.TrimSpace(part))
		if err != nil {
			fmt.Fprintf(os.Stderr, "parse error")
			return
		}
		numbers = append(numbers, num)
	}

	selected := selectFind(&numbers, n)
	fmt.Println("Selected number:", selected)
	fmt.Println("Number of swaps:", s)
	fmt.Println("Number of comparisons:", c)
}
