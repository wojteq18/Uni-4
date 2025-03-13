package queue

type Node struct {
	value int
	next  *Node
	prev  *Node
}

type CircularDll struct {
	head *Node
	size int
}

func (q *CircularDll) Push(value int) {
	newNode := &Node{value: value}
	if q.head == nil {
		newNode.next = newNode
		newNode.prev = newNode
		q.head = newNode
	} else {
		current := q.head
		for current.next != q.head {
			current = current.next
		}
		current.next = newNode
		newNode.prev = current
		newNode.next = q.head
		q.head.prev = newNode
	}
	q.size++
}

func (q *CircularDll) Remove() int {
	if q.head == nil {
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
		current.prev = q.head.prev
		q.head = q.head.next
	}
	q.size--
	return value
}

func Insert(list *CircularDll, value int) {
	list.Push(value)
}

func Merge(list1, list2 *CircularDll) *CircularDll {
	if list1.size == 0 {
		return list2
	} else if list2.size == 0 {
		return list1
	} else {
		// Connect the last element of the first list with the first element of the second list
		lastElement1 := list1.head.prev
		lastElement2 := list2.head.prev

		lastElement1.next = list2.head
		list2.head.prev = lastElement1

		lastElement2.next = list1.head
		list1.head.prev = lastElement2

		list1.size += list2.size
		return list1
	}

}

func Contains(list CircularDll, value int) (bool, int) {
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
