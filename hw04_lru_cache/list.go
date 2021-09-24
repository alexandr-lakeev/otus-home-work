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
	front *ListItem
	back  *ListItem
	len   int
}

func (l *list) Len() int {
	return l.len
}

func (l *list) Front() *ListItem {
	return l.front
}

func (l *list) Back() *ListItem {
	return l.back
}

func (l *list) PushFront(v interface{}) *ListItem {
	i := &ListItem{Value: v}

	if l.Len() == 0 {
		l.InitList(i)
		return i
	}

	l.PushFrontItem(i)
	l.len++

	return i
}

func (l *list) InitList(i *ListItem) {
	l.back, l.front = i, i
	l.len = 1
}

func (l *list) PushFrontItem(i *ListItem) {
	if l.Len() > 0 {
		oldFront := l.front
		i.Next, oldFront.Prev = oldFront, i
		l.front = i
	} else {
		l.InitList(i)
	}
}

func (l *list) PushBackItem(i *ListItem) {
	if l.Len() > 0 {
		oldBack := l.back
		i.Prev, oldBack.Next = oldBack, i
		l.back = i
	} else {
		l.InitList(i)
	}
}

func (l *list) PushBack(v interface{}) *ListItem {
	i := &ListItem{Value: v}

	if l.Len() == 0 {
		l.InitList(i)
		return i
	}

	l.PushBackItem(i)
	l.len++

	return i
}

func (l *list) Remove(i *ListItem) {
	if i == l.front {
		l.front = i.Next
	}

	if i == l.back {
		l.back = i.Prev
	}

	switch {
	case i.Next != nil && i.Prev != nil:
		i.Next.Prev, i.Prev.Next = i.Prev, i.Next
	case i.Next != nil:
		i.Next.Prev = nil
	case i.Prev != nil:
		i.Prev.Next = nil
	}

	i.Next, i.Prev = nil, nil
	l.len--
}

func (l *list) MoveToFront(i *ListItem) {
	l.Remove(i)
	l.PushFrontItem(i)
	l.len++
}

func NewList() List {
	return new(list)
}
