package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type BSTNode struct {
	Value  int
	Left   *BSTNode
	Right  *BSTNode
	Parent *BSTNode
}

type BST struct {
	Root *BSTNode
}

func (bst *BST) Height() int {
	if bst.Root == nil {
		return 0
	}

	type NodeLevel struct {
		Node  *BSTNode
		Level int
	}

	queue := []NodeLevel{{bst.Root, 0}}
	currentLevel := 0

	for len(queue) > 0 {
		nl := queue[0]
		queue = queue[1:]

		if nl.Level != currentLevel {
			currentLevel = nl.Level
		}

		if nl.Node.Left != nil {
			queue = append(queue, NodeLevel{nl.Node.Left, nl.Level + 1})
		}
		if nl.Node.Right != nil {
			queue = append(queue, NodeLevel{nl.Node.Right, nl.Level + 1})
		}
	}
	return currentLevel
}

func (bst *BST) Insert(value int) {
	fmt.Println("Insert ", value)
	x := bst.Root
	if x == nil {
		bst.Root = &BSTNode{Value: value}
		return
	}

	for x != nil {
		if value < x.Value {
			if x.Left == nil {
				x.Left = &BSTNode{Value: value, Parent: x}
				return
			}
			x = x.Left
		} else if value > x.Value {
			if x.Right == nil {
				x.Right = &BSTNode{Value: value, Parent: x}
				return
			}
			x = x.Right
		} else if value == x.Value {
			return
		}
	}
}

func (bst *BST) findSuccessor(node *BSTNode) *BSTNode {
	if node.Right != nil {
		node = node.Right
		for node.Left != nil {
			node = node.Left
		}
		return node
	}

	for node.Parent != nil && node == node.Parent.Right {
		node = node.Parent
	}
	return node.Parent
}

func (bst *BST) Search(value int) *BSTNode {
	x := bst.Root
	for x != nil {
		if value < x.Value {
			x = x.Left
		} else if value > x.Value {
			x = x.Right
		} else {
			return x
		}
	}
	return nil
}

func (bst *BST) Delete(value int) {
	fmt.Println("Delete ", value)
	if bst.Root == nil {
		return
	}

	node := bst.Search(value)
	if node == nil {
		return
	}

	//usuwany element jest rootem
	if bst.Root == node {
		if node.Left == nil && node.Right == nil {
			bst.Root = nil
			return
		}
		if node.Left != nil && node.Right == nil {
			bst.Root = node.Left
			node.Left.Parent = nil
			return
		}
		if node.Left == nil && node.Right != nil {
			bst.Root = node.Right
			node.Right.Parent = nil
			return
		}
		if node.Left != nil && node.Right != nil {
			successor := bst.findSuccessor(node)
			node.Value = successor.Value

			var child *BSTNode = nil
			if successor.Right != nil {
				child = successor.Right
			}

			if successor.Parent != nil {
				if successor == successor.Parent.Left {
					successor.Parent.Left = child
				} else {
					successor.Parent.Right = child
				}
			} else {
				bst.Root = child
			}

			if child != nil {
				child.Parent = successor.Parent
			}
			return
		}
	}

	//usuwany element jest liściem
	if node.Left == nil && node.Right == nil {
		if node.Parent != nil {
			if node == node.Parent.Left {
				node.Parent.Left = nil
			}
			if node == node.Parent.Right {
				node.Parent.Right = nil
			}
		}
		return
	}

	//usuwany element ma jedno dziecko
	if node.Left != nil && node.Right == nil {
		if node.Parent != nil {
			if node == node.Parent.Left {
				node.Parent.Left = node.Left
			}
			if node == node.Parent.Right {
				node.Parent.Right = node.Left
			}
		}
		node.Left.Parent = node.Parent
		return
	}

	if node.Left == nil && node.Right != nil {
		if node.Parent != nil {
			if node == node.Parent.Left {
				node.Parent.Left = node.Right
			}
			if node == node.Parent.Right {
				node.Parent.Right = node.Right
			}
		}
		node.Right.Parent = node.Parent
		return
	}

	//usuwany element ma dwoje dzieci
	if node.Left != nil && node.Right != nil {
		successor := bst.findSuccessor(node)
		node.Value = successor.Value

		var child *BSTNode = nil
		if successor.Right != nil {
			child = successor.Right
		}

		if successor.Parent != nil {
			if successor == successor.Parent.Left {
				successor.Parent.Left = child
			} else {
				successor.Parent.Right = child
			}
		} else {
			bst.Root = child
		}

		if child != nil {
			child.Parent = successor.Parent
		}
		return
	}
}

func (bst *BST) Print() {
	var printNode func(node *BSTNode, prefix string, isLeft bool)
	printNode = func(node *BSTNode, prefix string, isLeft bool) {
		if node == nil {
			return
		}

		printNode(node.Right, prefix+func() string {
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
		fmt.Printf("[%d]\n", node.Value)

		printNode(node.Left, prefix+func() string {
			if isLeft {
				return "|   "
			}
			return "    "
		}(), true)
	}

	printNode(bst.Root, "", false)
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

	bst := &BST{}

	// Insertuj
	for _, s := range insertStrs {
		num, err := strconv.Atoi(strings.TrimSpace(s))
		if err != nil {
			fmt.Printf("Błąd parsowania liczby do wstawienia: %v\n", s)
			continue
		}
		bst.Insert(num)
		bst.Print()
		fmt.Println()
	}

	// Usuwaj
	for _, s := range deleteStrs {
		num, err := strconv.Atoi(strings.TrimSpace(s))
		if err != nil {
			fmt.Printf("Błąd parsowania liczby do usunięcia: %v\n", s)
			continue
		}
		bst.Delete(num)
		bst.Print()
		fmt.Println()
	}
}
