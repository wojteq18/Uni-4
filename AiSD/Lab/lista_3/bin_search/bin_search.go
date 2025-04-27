package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func binarySearch(arr []int, target int) int {
	length := len(arr)

	if length == 0 {
		return -1
	}

	if length == 1 {
		if arr[0] == target {
			return 1
		} else {
			return -1
		}
	}

	pivot := arr[length/2]

	if pivot == target {
		return 1
	} else if target > pivot {
		return binarySearch(arr[length/2:], target)
	} else {
		return binarySearch(arr[:length/2], target)
	}
}

func main() {

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

	result := binarySearch(numbers, n)

	fmt.Println("Rezultat: ", result)

}
