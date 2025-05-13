package main

import (
	"fmt"
	"io"
	"math/rand"
	"os"
	"strconv"
	"strings"

	"github.com/seehuhn/mt19937"
)

var rng *rand.Rand

func init() {
	source := mt19937.New()
	source.Seed(123433334344)

	rng = rand.New(source)
}

func random_select(arr []int, n int, s *int, c *int) int {

	length := len(arr)

	if n > length {
		panic("n is greater than the length of the array")
	}

	randonPivot := rng.Intn(length)
	fmt.Println("Losowy pivot: ", (arr)[randonPivot])

	if randonPivot != length-1 {
		(arr)[randonPivot], (arr)[length-1] = (arr)[length-1], (arr)[randonPivot]
		(*s)++
	}

	i := 0

	for j := 0; j < length-1; j++ {
		(*c)++
		if (arr)[j] < (arr)[length-1] {
			if i != j {
				(*s)++
				(arr)[i], (arr)[j] = (arr)[j], (arr)[i]
			}
			i++
		}
	}

	if i != length-1 {
		(arr)[i], (arr)[length-1] = (arr)[length-1], (arr)[i]
		(*s)++
	}
	fmt.Println("Obecny stan: ", arr)

	if i < n-1 {
		right := (arr)[i+1:]
		return random_select(right, n-i-1, s, c)
	}

	if i > n-1 {
		left := (arr)[:i]
		return random_select(left, n, s, c)
	}

	return (arr)[i]

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
	fmt.Println("Tablica: ", numbers)

	i := random_select(numbers, n, &s, &c)

	fmt.Println("Znaleziony element: ", i)
	fmt.Println("s =", s)
	fmt.Println("c =", c)
}
