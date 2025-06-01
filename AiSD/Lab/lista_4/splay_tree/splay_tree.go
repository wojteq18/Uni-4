package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type STNode struct {
	value  int
	left   *STNode
	right  *STNode
	parent *STNode
}

type SplayTree struct {
	root *STNode
}

func (st *SplayTree) FindNode(value int) *STNode {
	current := st.root

	for current != nil {
		if value == current.value {
			return current
		} else if value < current.value {
			current = current.left
		} else {
			current = current.right
		}
	}
	return nil
}

func (st *SplayTree) leftRotate(node *STNode) {
	y := node.right
	node.right = y.left
	if y.left != nil {
		y.left.parent = node
	}
	y.parent = node.parent
	if node.parent == nil {
		st.root = y
	} else if node == node.parent.left {
		node.parent.left = y
	} else {
		node.parent.right = y
	}
	y.left = node
	node.parent = y
}

func (st *SplayTree) rightRotate(node *STNode) {
	y := node.left
	node.left = y.right
	if y.right != nil {
		y.right.parent = node
	}
	y.parent = node.parent
	if node.parent == nil {
		st.root = y
	} else if node == node.parent.right {
		node.parent.right = y
	} else {
		node.parent.left = y
	}
	y.right = node
	node.parent = y
}

func (st *SplayTree) Splay(value int) {
	node := st.FindNode(value)
	if node == nil {
		return
	}

	for node.parent != nil {
		if node.parent.parent == nil {
			// Zig
			if node == node.parent.left {
				st.rightRotate(node.parent)
			} else {
				st.leftRotate(node.parent)
			}
		} else if node == node.parent.left && node.parent == node.parent.parent.left {
			// Zig-Zig
			st.rightRotate(node.parent.parent)
			st.rightRotate(node.parent)
		} else if node == node.parent.right && node.parent == node.parent.parent.right {
			// Zig-Zig
			st.leftRotate(node.parent.parent)
			st.leftRotate(node.parent)
		} else if node == node.parent.right && node.parent == node.parent.parent.left {
			// Zig-Zag
			st.leftRotate(node.parent)
			st.rightRotate(node.parent)
		} else if node == node.parent.left && node.parent == node.parent.parent.right {
			// Zig-Zag
			st.rightRotate(node.parent)
			st.leftRotate(node.parent)
		}
	}
}

func (st *SplayTree) Insert(value int) {
	fmt.Println("Insert ", value)
	x := st.root
	if x == nil {
		st.root = &STNode{value: value}
		return
	}

	for x != nil {
		if value < x.value {
			if x.left == nil {
				x.left = &STNode{value: value, parent: x}
				st.Splay(value)
				return
			}
			x = x.left
		} else if value > x.value {
			if x.right == nil {
				x.right = &STNode{value: value, parent: x}
				st.Splay(value)
				return
			}
			x = x.right
		} else if value == x.value {
			st.Splay(value)
			return
		}
	}
	st.Splay(value)
}

func (st *SplayTree) Delete(value int) {
	fmt.Println("Delete ", value)
	node := st.FindNode(value)
	if node == nil {
		return
	}

	st.Splay(value)

	leftSubtree := node.left
	rightSubtree := node.right

	if leftSubtree != nil {
		leftSubtree.parent = nil
	}

	if rightSubtree != nil {
		rightSubtree.parent = nil
	}

	//Łączymy dwa poddrzewa
	st.Join(leftSubtree, rightSubtree)
}

func (st *SplayTree) Join(left, right *STNode) *STNode {
	if left == nil {
		return right
	}

	if right == nil {
		return left
	}

	maxNode := left
	for maxNode.right != nil {
		maxNode = maxNode.right
	}

	st.root = maxNode
	st.Splay(maxNode.value)

	st.root.right = right
	right.parent = st.root
	return st.root
}

func (st *SplayTree) Print() {
	var printNode func(node *STNode, prefix string, isLeft bool)
	printNode = func(node *STNode, prefix string, isLeft bool) {
		if node == nil {
			return
		}

		printNode(node.right, prefix+func() string {
			if isLeft {
				return "|   "
			}
			return "    "
		}(), false)

		fmt.Printf("%s", prefix)
		if isLeft {
			fmt.Printf("|-/")
		} else {
			fmt.Printf("\\-/")
		}
		fmt.Printf("[%d]\n", node.value)

		printNode(node.left, prefix+func() string {
			if isLeft {
				return "|   "
			}
			return "    "
		}(), true)
	}

	printNode(st.root, "", false)
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		fmt.Println("Brak danych wejściowych.")
		return
	}

	input := scanner.Text()
	parts := strings.SplitN(input, " ", 2)
	if len(parts) != 2 {
		fmt.Println("Podaj dane w formacie: insertowane_liczby usuwane_liczby")
		return
	}

	insertStrs := strings.Split(parts[0], ",")
	deleteStrs := strings.Split(parts[1], ",")

	st := &SplayTree{}

	// Insertuj
	for _, s := range insertStrs {
		num, err := strconv.Atoi(strings.TrimSpace(s))
		if err != nil {
			fmt.Printf("Błąd parsowania liczby do wstawienia: %v\n", s)
			continue
		}
		st.Insert(num)
		st.Print()
		fmt.Println()
	}

	// Usuwaj
	for _, s := range deleteStrs {
		num, err := strconv.Atoi(strings.TrimSpace(s))
		if err != nil {
			fmt.Printf("Błąd parsowania liczby do usunięcia: %v\n", s)
			continue
		}
		st.Delete(num)
		st.Print()
		fmt.Println()
	}
}
