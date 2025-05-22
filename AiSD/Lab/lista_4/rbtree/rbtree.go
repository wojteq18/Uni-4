package main

import (
	"fmt"
)

type RBNode struct {
	Value  int
	Color  rune
	Left   *RBNode
	Right  *RBNode
	Parent *RBNode
}

type RBTree struct {
	Root *RBNode
}

func (rbtree *RBTree) Search(value int) *RBNode {
	x := rbtree.Root
	for x != nil {
		if value == x.Value {
			return x
		} else if value < x.Value {
			x = x.Left
		} else {
			x = x.Right
		}
	}
	return nil
}

func (rbtree *RBTree) findSuccessor(node *RBNode) *RBNode {
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

func (rbtree *RBTree) Insert(value int) {
	//Pierwszy element - korzeń
	if rbtree.Root == nil {
		rbtree.Root = &RBNode{Value: value, Color: 'B'}
		return
	}

	//wstawianie nowego węzła
	x := rbtree.Root
	for x != nil {
		if value < x.Value {
			if x.Left == nil {
				x.Left = &RBNode{Value: value, Color: 'R', Parent: x}
				rbtree.fixInsert(x.Left)
				return
			}
			x = x.Left
		} else if value > x.Value {
			if x.Right == nil {
				x.Right = &RBNode{Value: value, Color: 'R', Parent: x}
				rbtree.fixInsert(x.Right)
				return
			}
			x = x.Right
		} else {
			return // Wartość już istnieje
		}
	}
}

func (rbtree *RBTree) fixInsert(node *RBNode) {
	//node jest rootem
	if node.Parent == nil {
		node.Color = 'B'
		return
	}

	// Rodzic jest czarny, nic nie trzeba robić
	if node.Parent.Color == 'B' {
		return
	}

	//Rodzic i Wujek są czerwoni
	if node.Parent.Parent != nil && node.Parent.Parent.Left != nil && node.Parent.Parent.Right != nil &&
		node.Parent.Parent.Left.Color == 'R' && node.Parent.Parent.Right.Color == 'R' {
		node.Parent.Parent.Right.Color = 'B'
		node.Parent.Parent.Left.Color = 'B'
		node.Parent.Parent.Color = 'R'
		rbtree.fixInsert(node.Parent.Parent)
		return
	}

	//Rodzic jest czerwony, Wujek jest czarny albo nil
	if node.Parent.Parent != nil && node.Parent.Color == 'R' && ((node.Parent.Parent.Left == nil || node.Parent.Parent.Left.Color != 'R') || (node.Parent.Parent.Right == nil || node.Parent.Parent.Right.Color != 'R')) {
		//wstawiany element jest wewnętrznym dzieckiem
		if node == node.Parent.Right && node.Parent == node.Parent.Parent.Left {
			rbtree.leftRotate((node.Parent))
			node = node.Left
		}
		if node == node.Parent.Left && node.Parent == node.Parent.Parent.Right {
			rbtree.rightRotate((node.Parent))
			node = node.Right
		}

		//wstawiany element jest zewnętrznym dzieckiem
		if node == node.Parent.Left && node.Parent.Parent.Left == node.Parent {
			node.Parent.Color = 'B'
			node.Parent.Parent.Color = 'R'
			rbtree.rightRotate(node.Parent.Parent)
		}
		if node == node.Parent.Right && node.Parent.Parent.Right == node.Parent {
			node.Parent.Color = 'B'
			node.Parent.Parent.Color = 'R'
			rbtree.leftRotate(node.Parent.Parent)
		}
	}
}

func (rbtree *RBTree) leftRotate(node *RBNode) {
	y := node.Right
	node.Right = y.Left
	if y.Left != nil {
		y.Left.Parent = node
	}
	y.Parent = node.Parent
	if node.Parent == nil {
		rbtree.Root = y
	} else if node == node.Parent.Left {
		node.Parent.Left = y
	} else {
		node.Parent.Right = y
	}
	y.Left = node
	node.Parent = y
}

func (rbtree *RBTree) rightRotate(node *RBNode) {
	y := node.Left
	node.Left = y.Right
	if y.Right != nil {
		y.Right.Parent = node
	}
	y.Parent = node.Parent
	if node.Parent == nil {
		rbtree.Root = y
	} else if node == node.Parent.Right {
		node.Parent.Right = y
	} else {
		node.Parent.Left = y
	}
	y.Right = node
	node.Parent = y
}

func (rbtree *RBTree) Delete(value int) {
	if rbtree.Root == nil {
		return
	}

	node := rbtree.Search(value)
	if node == nil {
		return
	}

	//usuwany element jest korzeniem
	if rbtree.Root == node {
		if node.Left == nil && node.Right == nil {
			rbtree.Root = nil
			return
		} else if node.Left == nil {
			rbtree.Root = node.Right
			rbtree.Root.Parent = nil
			if node.Color == 'B' {
				rbtree.fixDelete(rbtree.Root)
			}
			return
		} else if node.Right == nil {
			rbtree.Root = node.Left
			rbtree.Root.Parent = nil
			if node.Color == 'B' {
				rbtree.fixDelete(rbtree.Root)
			}
			return
		} else {
			succesor := rbtree.findSuccessor(node)
			originalColor := succesor.Color
			node.Value = succesor.Value

			var child *RBNode = nil
			if succesor.Right != nil {
				child = succesor.Right
			}

			if succesor.Parent != nil {
				if succesor == succesor.Parent.Left {
					succesor.Parent.Left = child
				} else {
					succesor.Parent.Right = child
				}
			} else {
				rbtree.Root = child
			}

			if child != nil {
				child.Parent = succesor.Parent
			}

			if originalColor == 'B' {
				rbtree.fixDelete(child)
			}
			return
		}
	}

	//usuwany element jest liściem
	if node.Left == nil && node.Right == nil {
		originalColor := node.Color
		if node.Parent != nil {
			if node == node.Parent.Left {
				node.Parent.Left = nil
			} else if node == node.Parent.Right {
				node.Parent.Right = nil
			}
		}
		if originalColor == 'B' {
			rbtree.fixDelete(node)
		}
		return
	}

	//usuwany element ma jedno dziecko
	if node.Left != nil && node.Right == nil {
		originalColor := node.Color
		if node.Parent != nil {
			if node == node.Parent.Left {
				node.Parent.Left = node.Left
			} else if node == node.Parent.Right {
				node.Parent.Right = node.Left
			}
		}
		node.Left.Parent = node.Parent
		if originalColor == 'B' {
			rbtree.fixDelete(node.Left)
		}
		return
	}
	if node.Left == nil && node.Right != nil {
		originalColor := node.Color
		if node.Parent != nil {
			if node == node.Parent.Left {
				node.Parent.Left = node.Right
			} else if node == node.Right.Parent {
				node.Parent.Right = node.Right
			}
		}
		node.Right.Parent = node.Parent
		if originalColor == 'B' {
			rbtree.fixDelete(node.Right)
		}
		return
	}

	//usuwany element ma dwoje dzieci
	if node.Left != nil && node.Right != nil {
		successor := rbtree.findSuccessor(node)
		originalColor := successor.Color
		node.Value = successor.Value

		var child *RBNode = nil
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
			rbtree.Root = child
		}

		if child != nil {
			child.Parent = successor.Parent
		}
		if originalColor == 'B' {
			rbtree.fixDelete(child)
		}
		return
	}
}

func (rbtree *RBTree) fixDelete(node *RBNode) {
	if node == nil {
		return
	}

	//jeśli node jest cerwony, zamień na czarny
	if node.Color == 'R' {
		node.Color = 'B'
		return
	}

	if node.Parent == nil {
		return
	}

	if node.Parent.Left == node {
		sibling := node.Parent.Right
		//Brat węzła x jest czerwony
		if sibling != nil && sibling.Color == 'R' {
			sibling.Color = 'B'
			node.Parent.Color = 'R'
			rbtree.leftRotate(node.Parent)
			sibling = node.Parent.Right
		}

		//Brat węzła x jest czarny i obydwoje dzieci są czarne
		if sibling != nil && sibling.Color == 'B' &&
			(sibling.Left == nil || sibling.Left.Color == 'B') &&
			(sibling.Right == nil || sibling.Right.Color == 'B') {
			sibling.Color = 'R'
			if node.Parent.Color == 'B' {
				rbtree.fixDelete(node.Parent)
			} else {
				node.Parent.Color = 'B'
			}
		} else {
			//Brat węzła x jest czarny, prawe dziecko brata jest czarne, lewe czerwone
			if sibling != nil && (sibling.Right == nil || sibling.Right.Color == 'B') {
				if sibling.Left != nil {
					sibling.Left.Color = 'B'
				}
				sibling.Color = 'R'
				rbtree.rightRotate(sibling)
				sibling = node.Parent.Right
			}

			//Brat węzła x jest czarny, prawe dziecko brata jest czerwone
			if sibling != nil {
				sibling.Color = node.Parent.Color
				node.Parent.Color = 'B'
				if sibling.Right != nil {
					sibling.Right.Color = 'B'
				}
				rbtree.leftRotate(node.Parent)
			}
		}
	} else {
		sibling := node.Parent.Left

		//brat węzła x jest czerwony
		if sibling != nil && sibling.Color == 'R' {
			sibling.Color = 'B'
			node.Parent.Color = 'R'
			rbtree.rightRotate(node.Parent)
			sibling = node.Parent.Left
		}

		//brat węzła x jest czarny i obydwoje dzieci są czarne
		if sibling != nil && sibling.Color == 'B' &&
			(sibling.Left == nil || sibling.Left.Color == 'B') &&
			(sibling.Right == nil || sibling.Right.Color == 'B') {
			sibling.Color = 'R'
			if node.Parent.Color == 'B' {
				rbtree.fixDelete(node.Parent)
			} else {
				node.Parent.Color = 'B'
			}
		} else {

			//brat węzła x jest czarny, lewe dziecko brata jest czarne, prawe czerwone
			if sibling != nil && (sibling.Left == nil || sibling.Left.Color == 'B') {
				if sibling.Right != nil {
					sibling.Right.Color = 'B'
				}
				sibling.Color = 'R'
				rbtree.leftRotate(sibling)
				sibling = node.Parent.Left
			}

			//brat węzła x jest czarny, lewe dziecko brata jest czerwone
			if sibling != nil {
				sibling.Color = node.Parent.Color
				node.Parent.Color = 'B'
				if sibling.Left != nil {
					sibling.Left.Color = 'B'
				}
				rbtree.rightRotate(node.Parent)
			}
		}
	}
}

func (rbtree *RBTree) Height() int {
	if rbtree.Root == nil {
		return 0
	}

	type NodeLevel struct {
		Node  *RBNode
		Level int
	}

	queue := []NodeLevel{{rbtree.Root, 0}}
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

func (bst *RBTree) Print() {
	var printNode func(node *RBNode, prefix string, isLeft bool)
	printNode = func(node *RBNode, prefix string, isLeft bool) {
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
		if node.Color == 'R' {
			fmt.Printf("\033[31m[%d]\033[0m\n", node.Value) // czerwony
		} else {
			fmt.Printf("[%d]\n", node.Value)
		}

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
	rbtree := &RBTree{}
	values := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	for _, value := range values {
		rbtree.Insert(value)
	}

	rbtree.Print()
	fmt.Println("Height of the tree:", rbtree.Height())
	rbtree.Delete(2)
	rbtree.Print()
}
