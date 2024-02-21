package hw04lrucache

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestListMethods(t *testing.T) {
	t.Run("push front", func(t *testing.T) {
		l := NewList()
		l.PushFront(10)
		l.PushFront(20)
		l.PushFront(30)
		require.Equal(t, 3, l.Len())

		elms := make([]int, 0, l.Len())
		for i := l.Front(); i != nil; i = i.Next {
			elms = append(elms, i.Value.(int))
		}
		require.Equal(t, elms, []int{30, 20, 10})

		elms = make([]int, 0, l.Len())
		for i := l.Back(); i != nil; i = i.Prev {
			elms = append(elms, i.Value.(int))
		}
		require.Equal(t, elms, []int{10, 20, 30})
	})

	t.Run("push back", func(t *testing.T) {
		l := NewList()
		l.PushBack(10)
		l.PushBack(20)
		l.PushBack(30)
		require.Equal(t, 3, l.Len())

		elms := make([]int, 0, l.Len())
		for i := l.Front(); i != nil; i = i.Next {
			elms = append(elms, i.Value.(int))
		}
		require.Equal(t, elms, []int{10, 20, 30})

		elms = make([]int, 0, l.Len())
		for i := l.Back(); i != nil; i = i.Prev {
			elms = append(elms, i.Value.(int))
		}
		require.Equal(t, elms, []int{30, 20, 10})
	})

	t.Run("remove once", func(t *testing.T) {
		l := NewList()
		l.Remove(l.PushFront(10))
		require.Equal(t, 0, l.Len())
	})

	t.Run("remove head", func(t *testing.T) {
		l := NewList()
		l.PushFront(10)
		l.Remove(l.PushFront(20))
		require.Equal(t, 1, l.Len())
		require.Nil(t, l.Front().Prev)
		require.Equal(t, l.Front().Value, 10)
	})

	t.Run("remove tail", func(t *testing.T) {
		l := NewList()
		l.PushFront(10)
		l.Remove(l.PushBack(20))
		require.Equal(t, 1, l.Len())
		require.Nil(t, l.Front().Next)
		require.Equal(t, l.Front().Value, 10)
	})

	t.Run("move once", func(t *testing.T) {
		l := NewList()
		l.MoveToFront(l.PushFront(10))
		require.Equal(t, 1, l.Len())
	})

	t.Run("move head", func(t *testing.T) {
		l := NewList()
		i := l.PushFront(10)
		l.PushBack(20)
		l.MoveToFront(i)
		require.Equal(t, 2, l.Len())
		require.Equal(t, l.Front().Next.Value, 20)
	})

	t.Run("move tail", func(t *testing.T) {
		l := NewList()
		i := l.PushFront(10)
		l.PushFront(20)
		l.MoveToFront(i)
		require.Equal(t, 2, l.Len())
		require.Equal(t, l.Front().Next.Value, 20)
	})

	t.Run("move middle", func(t *testing.T) {
		l := NewList()
		for _, v := range []int{10, 20, 30, 40, 50} {
			l.PushBack(v)
		}
		l.MoveToFront(l.Front().Next.Next) // 30
		elms := make([]int, 0, l.Len())
		for i := l.Front(); i != nil; i = i.Next {
			elms = append(elms, i.Value.(int))
		}
		require.Equal(t, elms, []int{30, 10, 20, 40, 50})

		elms = make([]int, 0, l.Len())
		for i := l.Back(); i != nil; i = i.Prev {
			elms = append(elms, i.Value.(int))
		}
		require.Equal(t, elms, []int{50, 40, 20, 10, 30})
	})
}

func TestList(t *testing.T) {
	t.Run("empty list", func(t *testing.T) {
		l := NewList()

		require.Equal(t, 0, l.Len())
		require.Nil(t, l.Front())
		require.Nil(t, l.Back())
	})

	t.Run("complex", func(t *testing.T) {
		l := NewList()

		l.PushFront(10) // [10]
		l.PushBack(20)  // [10, 20]
		l.PushBack(30)  // [10, 20, 30]
		require.Equal(t, 3, l.Len())

		middle := l.Front().Next // 20
		l.Remove(middle)         // [10, 30]
		require.Equal(t, 2, l.Len())

		for i, v := range [...]int{40, 50, 60, 70, 80} {
			if i%2 == 0 {
				l.PushFront(v)
			} else {
				l.PushBack(v)
			}
		} // [80, 60, 40, 10, 30, 50, 70]

		require.Equal(t, 7, l.Len())
		require.Equal(t, 80, l.Front().Value)
		require.Equal(t, 70, l.Back().Value)

		l.MoveToFront(l.Front()) // [80, 60, 40, 10, 30, 50, 70]
		l.MoveToFront(l.Back())  // [70, 80, 60, 40, 10, 30, 50]

		elems := make([]int, 0, l.Len())
		for i := l.Front(); i != nil; i = i.Next {
			elems = append(elems, i.Value.(int))
		}
		require.Equal(t, []int{70, 80, 60, 40, 10, 30, 50}, elems)
	})
}
