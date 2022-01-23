package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	firstItem *ListItem
}

func (list list) Len() int {
	count := 0
	item := list.Front()
	for item != nil {
		count++
		item = item.Next
	}
	return count
}

func (list list) Front() *ListItem {
	itemValue := list.firstItem
	if itemValue == nil {
		return itemValue
	}
	for itemValue.Prev != nil {
		itemValue = itemValue.Prev
	}

	return itemValue
}

func (list list) Back() *ListItem {
	itemValue := list.firstItem
	if itemValue == nil {
		return itemValue
	}
	for itemValue.Next != nil {
		itemValue = itemValue.Next
	}

	return itemValue
}

func (list *list) PushFront(v interface{}) *ListItem {
	firstItem := list.Front()
	if firstItem != nil {
		currentPrev := &ListItem{
			Value: v,
		}
		list.linkValues(currentPrev, firstItem)

		return currentPrev

	}
	list.firstItem = &ListItem{
		Value: v,
	}
	return list.firstItem
}

func (list *list) PushBack(v interface{}) *ListItem {
	nextValue := &ListItem{
		Value: v,
	}
	list.linkValues(list.Back(), nextValue)

	return nextValue
}

func (list *list) Remove(item *ListItem) {
	list.linkValues(item.Prev, item.Next)
	if item == list.firstItem {
		if item.Prev != nil {
			list.firstItem = item.Prev
		} else {

			list.firstItem = item.Next
		}
	}
	item.Next = nil
	item.Prev = nil
}

func (list *list) MoveToFront(item *ListItem) {
	if list.Len() > 1 {
		list.Remove(item)
		item.Prev = nil
		list.linkValues(item, list.Front())
	}
}

func NewList() List {
	return new(list)
}

func (list *list) linkValues(prevValue, nextValue *ListItem) *ListItem {
	if prevValue != nil {
		prevValue.Next = nextValue
	}
	if nextValue != nil {
		nextValue.Prev = prevValue
	}
	//При изменения порядка, изменяем первый элемент
	list.firstItem = list.Front()

	return prevValue
}
