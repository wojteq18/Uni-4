package queue

import "testing"

func TestInsert(t *testing.T) {
	list := CircularDll{}
	Insert(&list, 1)
	Insert(&list, 2)
	Insert(&list, 5)
	if list.size != 3 {
		t.Errorf("Expected size: %d, got: %d", 3, list.size)
	}
	if list.head.prev.value != 5 {
		t.Errorf("Expected value: %d, got: %d", 3, list.head.value)
	}
}

func TestRemove(t *testing.T) {
	list := CircularDll{}
	Insert(&list, 1)
	Insert(&list, 2)
	Insert(&list, 5)
	removed := list.Remove()
	if removed != 1 {
		t.Errorf("Expected removed value: %d, got: %d", 5, removed)
	}
	if list.size != 2 {
		t.Errorf("Expected size: %d, got: %d", 2, list.size)
	}
	if list.head.value != 2 {
		t.Errorf("Expected head value: %d, got: %d", 2, list.head.value)
	}
}

func TestMerge(t *testing.T) {
	list1 := CircularDll{}
	list2 := CircularDll{}
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
	if list3.head.prev.value != 50 {
		t.Errorf("Expected head value: %d, got: %d", 50, list3.head.value)
	}
}

func TestContains(t *testing.T) {
	list := CircularDll{}
	Insert(&list, 1)
	Insert(&list, 2)
	Insert(&list, 5)

	found_2, _ := Contains(list, 2)
	if found_2 == false {
		t.Errorf("Expected true, got false")
	}

	found_3, _ := Contains(list, 3)
	if found_3 == true {
		t.Errorf("Expected false, got true")
	}
}
