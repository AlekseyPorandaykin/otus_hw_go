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

func NewListItem(v interface{}) *ListItem {
	return &ListItem{
		Value: v,
	}
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
	if list.isEmptyList() {
		return list.addFirstValue(v)
	}
	firstItem := NewListItem(v)
	list.linkValues(firstItem, list.Front())

	return firstItem
}

func (list *list) PushBack(v interface{}) *ListItem {
	list.count++
	if list.isEmptyList() {
		return list.addFirstValue(v)
	}
	lastItem := NewListItem(v)
	list.linkValues(list.Back(), lastItem)

	return lastItem
}

func (list *list) Remove(item *ListItem) {
	if list.count == 1 && list.firstItem == item {
		list.clear()
		return
	}
	if item == list.firstItem {
		list.firstItem = item.Next
	}
	if item == list.lastItem {
		list.lastItem = item.Prev
	}
	list.linkValues(item.Prev, item.Next)

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

func (list *list) addFirstValue(v interface{}) *ListItem {
	list.firstItem = NewListItem(v)
	list.lastItem = list.firstItem
	return list.firstItem
}

func (list *list) clear() {
	list.firstItem = nil
	list.lastItem = nil
	list.count = 0
}

func (list list) isEmptyList() bool {
	return list.Back() == nil && list.Front() == nil
}
