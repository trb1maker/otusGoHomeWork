package hw04lrucache

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	t.Run("empty cache", func(t *testing.T) {
		c := NewCache(10)

		_, ok := c.Get("aaa")
		require.False(t, ok)

		_, ok = c.Get("bbb")
		require.False(t, ok)
	})

	t.Run("simple", func(t *testing.T) {
		c := NewCache(5)

		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)

		wasInCache = c.Set("bbb", 200)
		require.False(t, wasInCache)

		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 100, val)

		val, ok = c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, 200, val)

		wasInCache = c.Set("aaa", 300)
		require.True(t, wasInCache)

		val, ok = c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 300, val)

		val, ok = c.Get("ccc")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("purge logic", func(t *testing.T) {
		c := NewCache(3)
		c.Set("a", 10) // [a]
		c.Set("b", 20) // [b, a]
		c.Set("c", 30) // [c, b, a]
		c.Get("a")     // [a, c, b]
		c.Set("d", 40) // [d, a, c]

		v, ok := c.Get("b")
		require.Nil(t, v)
		require.False(t, ok)

		c.Get("a")     // [a, d, c]
		c.Set("e", 20) // [e, a, d]

		v, ok = c.Get("c")
		require.Nil(t, v)
		require.False(t, ok)

		c.Clear()
		v, ok = c.Get("a")
		require.Nil(t, v)
		require.False(t, ok)

		c.Set("a", 10) // [a]
		c.Set("b", 20) // [b, a]
		c.Set("c", 30) // [c, b, a]
		c.Get("a")     // [a, c, b]
		c.Set("d", 40) // [d, a, c]

		v, ok = c.Get("b")
		require.Nil(t, v)
		require.False(t, ok)
	})
}

func TestCacheMultithreading(t *testing.T) {
	t.Skip() // Remove me if task with asterisk completed.

	c := NewCache(10)
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Set(Key(strconv.Itoa(i)), i)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Get(Key(strconv.Itoa(rand.Intn(1_000_000))))
		}
	}()

	wg.Wait()
}
