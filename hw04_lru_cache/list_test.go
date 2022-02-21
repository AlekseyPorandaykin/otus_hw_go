package hw04lrucache

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestList(t *testing.T) {
	t.Run("empty list", func(t *testing.T) {
		l := NewList()
		l.PushBack(45)
		tr := l.Front()
		l.Remove(tr)

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

	t.Run("testing boundary elements", func(t *testing.T) {
		l := NewList()
		singularValue := l.PushBack("singular value")
		l.Remove(singularValue)
		require.Empty(t, l.Len(), "List must be empty")
		firstItem := l.PushBack(1)
		secondItem := l.PushBack(2)
		l.PushBack(3)
		penultimateItem := l.PushBack(4)
		lastItem := l.PushBack(5)
		l.Remove(firstItem)
		require.True(
			t,
			l.Front() == secondItem,
			"When deleting the first value, the next value must become the first",
		)
		l.Remove(lastItem)
		require.True(
			t,
			l.Back() == penultimateItem,
			"When the last value is removed, the last element is set to the penultimate one.",
		)
	})
}
