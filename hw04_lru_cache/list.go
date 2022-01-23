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
	lastItem  *ListItem
	count     int
}

func (list list) Len() int {
	return list.count
}

func (list list) Front() *ListItem {
	return list.firstItem
}

func (list list) Back() *ListItem {
	return list.lastItem
}

func (list *list) PushFront(v interface{}) *ListItem {
	list.count++
	if list.Front() != nil {
		currentPrev := &ListItem{
			Value: v,
		}
		list.linkValues(currentPrev, list.Front())

		return currentPrev
	}
	list.firstItem = &ListItem{
		Value: v,
	}
	list.lastItem = list.firstItem
	return list.firstItem
}

func (list *list) PushBack(v interface{}) *ListItem {
	nextValue := &ListItem{
		Value: v,
	}
	list.linkValues(list.Back(), nextValue)
	list.count++

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

	list.count--
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

func (list *list) linkValues(prevValue, nextValue *ListItem) {
	if prevValue != nil {
		prevValue.Next = nextValue
		if prevValue.Prev == nil {
			list.firstItem = prevValue
		}
	}
	if nextValue != nil {
		nextValue.Prev = prevValue

		if nextValue.Next == nil {
			list.lastItem = nextValue
		}
	}
}
