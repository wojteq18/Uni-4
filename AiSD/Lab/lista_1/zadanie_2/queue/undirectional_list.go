package queue

type Node struct { // Node is a struct that contains an integer value and a pointer to the next node
	value int
	next  *Node
}

type UndirectionalList struct { // UndirectionalList is a struct that contains a pointer to the head node and the size of the list
	head *Node
	size int
}

func (q *UndirectionalList) Push(value int) {
	newNode := &Node{value: value} // Create a new node
	if q.head == nil {
		q.head = newNode
		newNode.next = newNode // Point the next node to itself
	} else {
		current := q.head
		for current.next != q.head {
			current = current.next
		}
		newNode.next = q.head
		current.next = newNode
		q.head = newNode
	}
	q.size++
}

func (q *UndirectionalList) Remove() int {
	if q.size == 0 {
		panic("The list is empty")
	}
	value := q.head.value
	if q.size == 1 {
		q.head = nil
	} else {
		current := q.head
		for current.next != q.head {
			current = current.next
		}
		current.next = q.head.next
		q.head = q.head.next
	}
	q.size--
	return value
}

func Insert(list *UndirectionalList, value int) {
	list.Push(value)
}

func Merge(list1, list2 UndirectionalList) UndirectionalList {
	if list1.size == 0 {
		return list2
	} else if list2.size == 0 {
		return list1
	} else {
		current := list2.head
		for current.next != list2.head {
			Insert(&list1, current.value)
			current = current.next
		}
		return list1
	}
}
