package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func sortFives(arr *[]int, s *int, c *int) *[]int {
	length := len(*arr)
	for i := 0; i < length; i++ {
		j := i
		for {
			//(*c)++
			if j <= 0 {
				break
			}
			//(*c)++
			if (*arr)[j] >= (*arr)[j-1] {
				break
			}
			(*arr)[j], (*arr)[j-1] = (*arr)[j-1], (*arr)[j]
			//(*s)++
			j--
		}
	}
	return arr
}

func findMedianInFive(arr *[]int, s *int, c *int) int {
	sortFives(arr, s, c)
	length := len(*arr)

	if length == 0 {
		panic("Array is empty")
	}

	return (*arr)[length/2]
}

func medianOfMedians(arr *[]int, s *int, c *int) int {

	length := len(*arr)
	copyArr := make([]int, length)
	copy(copyArr, *arr)

	if length <= 5 {
		return findMedianInFive(&copyArr, s, c)
	}

	var medians []int
	for i := 0; i < length; i += 5 {
		end := i + 5
		if end > length {
			end = length
		}
		subArr := make([]int, end-i)
		copy(subArr, (*arr)[i:end])
		medians = append(medians, findMedianInFive(&subArr, s, c))
	}

	return medianOfMedians(&medians, s, c)
}

func selectFind(arr *[]int, k int, s *int, c *int) int {
	length := len(*arr)
	pivotValue := medianOfMedians(arr, s, c)

	idx := 0

	for i := 0; i < length; i++ {
		(*c)++
		if (*arr)[i] == pivotValue {
			idx = i
			break
		}
	}

	if idx != length-1 {
		(*arr)[idx], (*arr)[length-1] = (*arr)[length-1], (*arr)[idx]
		(*s)++
	}

	i := 0
	for j := 0; j < length-1; j++ {
		(*c)++
		if (*arr)[j] < (*arr)[length-1] {
			if i != j {
				(*arr)[i], (*arr)[j] = (*arr)[j], (*arr)[i]
				(*s)++
			}
			i++
		}
	}

	if i != length-1 {
		(*arr)[i], (*arr)[length-1] = (*arr)[length-1], (*arr)[i]
		(*s)++
	}

	if i < k-1 {
		right := (*arr)[i+1:]
		return selectFind(&right, k-i-1, s, c)
	}

	if i > k-1 {
		left := (*arr)[:i]
		return selectFind(&left, k, s, c)
	}

	return (*arr)[i]
}

func main() {
	s := 0
	c := 0

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

	a := selectFind(&numbers, n, &s, &c)

	fmt.Println("s =", s)
	fmt.Println("c =", c)
	fmt.Println("znaleziony element:", a)
}

//../random_generator/target/release/random_generator 11 | go run random_select.go
