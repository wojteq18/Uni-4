package main

import (
	"fmt"
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

func (bst *BST) printLevels() {
	if bst.Root == nil {
		fmt.Println("Tree is empty")
		return
	}

	type NodeLevel struct {
		Node  *BSTNode
		Level int
	}

	queue := []NodeLevel{{bst.Root, 0}}
	currentLevel := 0

	fmt.Printf("Level %d: ", currentLevel)

	for len(queue) > 0 {
		nl := queue[0]
		queue = queue[1:]

		if nl.Level != currentLevel {
			currentLevel = nl.Level
			fmt.Printf("\nLevel %d: ", currentLevel)
		}

		fmt.Printf("%d ", nl.Node.Value)

		if nl.Node.Left != nil {
			queue = append(queue, NodeLevel{nl.Node.Left, nl.Level + 1})
		}
		if nl.Node.Right != nil {
			queue = append(queue, NodeLevel{nl.Node.Right, nl.Level + 1})
		}
	}

	fmt.Println()
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
		}
		if value > x.Value {
			if x.Right == nil {
				x.Right = &BSTNode{Value: value, Parent: x}
				return
			}
			x = x.Right
		}
		if value == x.Value {
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

	if bst.Root == nil {
		return
	}

	node := bst.Search(value)
	if node.Right == nil && node.Left == nil {
		Parent := node.Parent
		if Parent.Left == node {
			Parent.Left = nil
		} else {
			Parent.Right = nil
		}
	}

}

func main() {
	bst := &BST{}
	bst.Insert(10)
	bst.Insert(5)
	bst.Insert(15)
	bst.Insert(3)
	bst.Insert(2)
	bst.Insert(1)
	bst.Insert(4)

	bst.printLevels()
	bst.Delete(1)
	bst.printLevels()
}
