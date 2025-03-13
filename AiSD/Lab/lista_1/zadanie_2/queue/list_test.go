package queue

import "testing"

func testPush(t *testing.T) {

	list := UndirectionalList{}
	list.Push(1)
	list.Push(2)
	list.Push(5)

	if list.size != 3 {
		t.Errorf("Expected size: %d, got: %d", 3, list.size)
	}
	if list.head.value != 1 {
		t.Errorf("Expected value: %d, got: %d", 1, list.head.value)
	}
	if list.head.next.value != 2 {
		t.Errorf("Expected value: %d, got: %d", 2, list.head.next.value)
	}
}

func testRemove(t *testing.T) {
	list := UndirectionalList{}
	list.Push(1)
	list.Push(2)
	list.Push(5)

	removed1 := list.Remove()
	removed2 := list.Remove()

	if removed1 != 1 {
		t.Errorf("Expected removed value: %d, got: %d", 1, removed1)
	}
	if removed2 != 2 {
		t.Errorf("Expected removed value: %d, got: %d", 2, removed2)
	}
}

func TestMerge(t *testing.T) {
	list1 := UndirectionalList{}
	list2 := UndirectionalList{}
	Insert(&list1, 1)
	Insert(&list1, 2)
	Insert(&list1, 5)
	Insert(&list2, 10)
	Insert(&list2, 20)
	Insert(&list2, 50)
	list3 := Merge(&list1, &list2)

	if list3.size != 6 {
		t.Errorf("Expected size: %d, got: %d", 6, list3.size)
	}
	if list3.head.value != 1 {
		t.Errorf("Expected head value: %d, got: %d", 1, list3.head.value)
	}
}

func TestContains(t *testing.T) {
	list := UndirectionalList{}
	Insert(&list, 1)
	Insert(&list, 2)
	Insert(&list, 5)

	found_1, _ := Contains(list, 1)
	if found_1 == false {
		t.Errorf("Expected true, got false")
	}

	found_3, _ := Contains(list, 3)
	if found_3 == true {
		t.Errorf("Expected false, got true")
	}
}
