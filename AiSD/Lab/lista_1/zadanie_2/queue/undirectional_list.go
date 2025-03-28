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
		for current.next != q.head { // Find the last element
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

func Merge(list1, list2 *UndirectionalList) *UndirectionalList {
	if list1.size == 0 {
		return list2
	} else if list2.size == 0 {
		return list1
	} else {
		lastElement1 := list1.head
		for lastElement1.next != list1.head { // Find the last element of the first list
			lastElement1 = lastElement1.next
		}

		lastElement2 := list2.head
		for lastElement2.next != list2.head { // Find the last element of the second list
			lastElement2 = lastElement2.next
		}

		// Merge the lists
		lastElement1.next = list2.head
		lastElement2.next = list1.head

		list1.size += list2.size
		return list1
	}
}

func Contains(list UndirectionalList, value int) (bool, int) {
	comparsion := 0
	if list.head == nil {
		return false, 1
	}
	current := list.head
	for { // Check if the list contains the value
		comparsion++
		if current.value == value {
			return true, comparsion
		}
		current = current.next
		if current == list.head {
			break
		}
	}
	return false, comparsion
}
