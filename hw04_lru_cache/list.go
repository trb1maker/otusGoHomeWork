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
	i.Next, l.head.Prev = l.head, i
	l.head = i
	l.length++
	return i
}

func (l *list) PushBack(v interface{}) *ListItem {
	i := &ListItem{Value: v}
	if l.length == 0 {
		return l.pushFirst(i)
	}
	i.Prev, l.tail.Next = l.tail, i
	l.tail = i
	l.length++
	return i
}

func (l *list) clear() {
	l.head, l.tail = nil, nil
	l.length = 0
}

func (l *list) removeHead() {
	l.head.Next.Prev, l.head = nil, l.head.Next
	l.length--
}

func (l *list) removeTail() {
	l.tail.Prev.Next, l.tail = nil, l.tail.Prev
	l.length--
}

func (l *list) Remove(i *ListItem) {
	if l.length == 1 {
		l.clear()
		return
	}
	if l.head == i {
		l.removeHead()
		return
	}
	if l.tail == i {
		l.removeTail()
		return
	}
	i.Prev.Next, i.Next.Prev = i.Next, i.Prev
	l.length--
}

func (l *list) MoveToFront(i *ListItem) {
	if l.length == 1 || l.head == i {
		return
	}
	if l.tail == i {
		tail := i.Prev
		i.Next, l.head.Prev = l.head, i
		i.Prev, tail.Next = nil, nil
		l.head, l.tail = i, tail
		return
	}
	i.Prev.Next, i.Next.Prev = i.Next, i.Prev
	i.Prev, i.Next = nil, l.head
	l.head.Prev, l.head = i, i
}

func NewList() List {
	return new(list)
}
