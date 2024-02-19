package hw04lrucache

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
	keys     map[*ListItem]Key
}

func (l *lruCache) Set(key Key, value interface{}) bool {
	if i, ok := l.items[key]; ok {
		l.queue.MoveToFront(i)
		i.Value = value
		return true
	}
	if l.queue.Len() == l.capacity {
		unused := l.queue.Back()
		l.queue.Remove(unused)
		delete(l.items, l.keys[unused])
	}
	item := l.queue.PushFront(value)
	l.items[key], l.keys[item] = item, key
	return false
}

func (l *lruCache) Get(key Key) (interface{}, bool) {
	if i, ok := l.items[key]; ok {
		l.queue.MoveToFront(i)
		return i.Value, true
	}
	return nil, false
}

func (l *lruCache) Clear() {
	l.items = make(map[Key]*ListItem, l.capacity)
	l.queue = NewList()
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
		keys:     make(map[*ListItem]Key),
	}
}
