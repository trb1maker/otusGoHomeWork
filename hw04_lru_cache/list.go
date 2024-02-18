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
	head   *ListItem
	tail   *ListItem
	length int
}

func (l *list) Len() int {
	return l.length
}

func (l *list) Front() *ListItem {
	return l.head
}

func (l *list) Back() *ListItem {
	return l.tail
}

func (l *list) pushFirst(i *ListItem) *ListItem {
	l.head, l.tail = i, i
	l.length++
	return i
}

func (l *list) PushFront(v interface{}) *ListItem {
	i := &ListItem{Value: v}
	if l.length == 0 {
		return l.pushFirst(i)
	}
	i.Next = l.head
	l.head.Prev = i
	l.head = i
	l.length++
	return i
}

func (l *list) PushBack(v interface{}) *ListItem {
	i := &ListItem{Value: v}
	if l.length == 0 {
		return l.pushFirst(i)
	}
	i.Prev = l.tail
	l.tail.Next = i
	l.tail = i
	l.length++
	return i
}

func (l *list) Remove(i *ListItem) {
	if l.length == 1 {
		l.head = nil
		l.tail = nil
		l.length--
		return
	}
	if l.head == i {
		l.head.Next.Prev = nil
		l.head = l.head.Next
		l.length--
		return
	}
	if l.tail == i {
		l.tail.Prev.Next = nil
		l.tail = l.tail.Prev
		l.length--
		return
	}
	prev, next := i.Prev, i.Next
	prev.Next = next
	next.Prev = prev
	l.length--
}
func (l *list) MoveToFront(i *ListItem) {
	if l.length == 1 {
		return
	}
	l.Remove(i)
	l.PushFront(i.Value)
}

func NewList() List {
	return new(list)
}
