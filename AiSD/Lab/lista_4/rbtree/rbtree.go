package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
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
	fmt.Println("Insert ", value)
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

// Pomocnicza funkcja do pobierania koloru węzła (nil jest traktowany jako czarny)
func getNodeColor(node *RBNode) rune {
	if node == nil {
		return 'B'
	}
	return node.Color
}

func (rbtree *RBTree) Delete(value int) {
	fmt.Println("Delete ", value)
	nodeToDelete := rbtree.Search(value)
	if nodeToDelete == nil {
		return // Węzeł nie znaleziony
	}

	var y *RBNode // Węzeł, który faktycznie zostanie usunięty lub przesunięty
	var x *RBNode // Dziecko węzła y, które zajmie jego miejsce
	var yOriginalColor rune

	// Krok 1: Znajdź węzeł 'y', który zostanie usunięty ze struktury drzewa.
	// Jeśli nodeToDelete ma dwoje dzieci, 'y' będzie jego następnikiem.
	// W przeciwnym razie 'y' to sam nodeToDelete.
	if nodeToDelete.Left == nil || nodeToDelete.Right == nil {
		y = nodeToDelete
	} else {
		y = rbtree.findSuccessor(nodeToDelete)
	}
	yOriginalColor = y.Color

	// Krok 2: 'x' to jedyne dziecko 'y' (lub nil, jeśli 'y' nie ma dzieci).
	if y.Left != nil {
		x = y.Left
	} else {
		x = y.Right
	}

	// xParent to rodzic 'y'. Będzie on rodzicem 'x' po usunięciu 'y'.
	// Jest to kluczowe dla fixDelete, zwłaszcza gdy x jest nil.
	xParent := y.Parent

	// Krok 3: Usuń 'y' ze struktury drzewa, zastępując go przez 'x'.
	if x != nil {
		x.Parent = y.Parent
	}

	if y.Parent == nil { // Jeśli 'y' był korzeniem
		rbtree.Root = x
	} else { // Jeśli 'y' nie był korzeniem
		if y == y.Parent.Left {
			y.Parent.Left = x
		} else {
			y.Parent.Right = x
		}
	}

	// Krok 4: Jeśli 'y' był następnikiem nodeToDelete, skopiuj dane z 'y' do nodeToDelete.
	if y != nodeToDelete {
		nodeToDelete.Value = y.Value
		// Kolor nodeToDelete nie zmienia się tutaj, tylko jego wartość.
		// Balansowanie zależy od oryginalnego koloru 'y'.
	}

	// Krok 5: Jeśli usunięty węzeł 'y' był czarny, mogło dojść do naruszenia własności drzewa.
	if yOriginalColor == 'B' {
		// 'x' to węzeł, który zajął miejsce 'y'.
		// 'xParent' to rodzic nowej pozycji 'x' (były rodzic 'y').
		rbtree.fixDelete(x, xParent)
	}
}

func (rbtree *RBTree) fixDelete(x *RBNode, xParent *RBNode) {
	var sibling *RBNode

	for x != rbtree.Root && getNodeColor(x) == 'B' {
		if xParent == nil {
			break
		}

		if x == xParent.Left { // 'x' jest lewym dzieckiem (lub nil w miejscu lewego dziecka)
			sibling = xParent.Right

			if getNodeColor(sibling) == 'R' { // Przypadek 1: Brat jest czerwony
				if sibling != nil {
					sibling.Color = 'B'
				}
				xParent.Color = 'R'
				rbtree.leftRotate(xParent)
				sibling = xParent.Right // Brat się zmienił, zaktualizuj
			}

			if (sibling == nil || getNodeColor(sibling.Left) == 'B') &&
				(sibling == nil || getNodeColor(sibling.Right) == 'B') {
				if sibling != nil {
					sibling.Color = 'R'
				}
				x = xParent // Przesuń problem w górę
				if x != nil {
					xParent = x.Parent
				} else {
					xParent = nil
				} // Zaktualizuj xParent dla nowego x
			} else { // Brat jest czarny, przynajmniej jedno dziecko czerwone
				// Przypadek 3: Brat jest czarny, jego prawe dziecko jest czarne (lewe dziecko musi być czerwone)
				if sibling != nil && getNodeColor(sibling.Right) == 'B' { // sibling.Left musi być czerwone
					if sibling.Left != nil {
						sibling.Left.Color = 'B'
					}
					sibling.Color = 'R'
					rbtree.rightRotate(sibling)
					sibling = xParent.Right // Brat się zmienił
				}
				// Przypadek 4: Brat jest czarny, jego prawe dziecko jest czerwone
				if sibling != nil {
					sibling.Color = getNodeColor(xParent)
				}
				xParent.Color = 'B'
				if sibling != nil && sibling.Right != nil {
					sibling.Right.Color = 'B'
				}
				rbtree.leftRotate(xParent)
				x = rbtree.Root // Aby zakończyć pętlę
			}
		} else { // Symetryczny przypadek: x == xParent.Right (lub nil i był prawym dzieckiem)
			sibling = xParent.Left

			if getNodeColor(sibling) == 'R' { // Przypadek 1
				if sibling != nil {
					sibling.Color = 'B'
				}
				xParent.Color = 'R'
				rbtree.rightRotate(xParent)
				sibling = xParent.Left
			}

			if (sibling == nil || getNodeColor(sibling.Left) == 'B') &&
				(sibling == nil || getNodeColor(sibling.Right) == 'B') { // Przypadek 2
				if sibling != nil {
					sibling.Color = 'R'
				}
				x = xParent
				if x != nil {
					xParent = x.Parent
				} else {
					xParent = nil
				}
			} else {
				// Przypadek 3: Brat jest czarny, jego lewe dziecko jest czarne (prawe dziecko musi być czerwone)
				if sibling != nil && getNodeColor(sibling.Left) == 'B' { // sibling.Right musi być czerwone
					if sibling.Right != nil {
						sibling.Right.Color = 'B'
					}
					sibling.Color = 'R'
					rbtree.leftRotate(sibling)
					sibling = xParent.Left
				}
				// Przypadek 4: Brat jest czarny, jego lewe dziecko jest czerwone
				if sibling != nil {
					sibling.Color = getNodeColor(xParent)
				}
				xParent.Color = 'B'
				if sibling != nil && sibling.Left != nil {
					sibling.Left.Color = 'B'
				}
				rbtree.rightRotate(xParent)
				x = rbtree.Root
			}
		}
	}
	if x != nil {
		x.Color = 'B' //korzeń musi być czarny
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

	rbtree := &RBTree{}

	for _, s := range insertStrs {
		num, err := strconv.Atoi(strings.TrimSpace(s))
		if err != nil {
			fmt.Printf("Błąd parsowania liczby do wstawienia: %v\n", s)
			continue
		}
		rbtree.Insert(num)
		rbtree.Print()
		fmt.Println()
	}

	for _, s := range deleteStrs {
		num, err := strconv.Atoi(strings.TrimSpace(s))
		if err != nil {
			fmt.Printf("Błąd parsowania liczby do usunięcia: %v\n", s)
			continue
		}
		rbtree.Delete(num)
		rbtree.Print()
		fmt.Println()
	}
}
