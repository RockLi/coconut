package lfu

import (
	"container/list"
	"github.com/flatpeach/coconut/cache"
	"sync"
)

// This package implements a O(1) LFU
// Followed the origin paper http://dhruvbird.com/lfu.pdf

type Cache struct {
	mu sync.Mutex

	capacity    uint64
	size        uint64
	maxElements uint64
	freq        *list.List
	caches      map[cache.Key]*entry
}

type node struct {
	freq  int
	items map[cache.Key]uint8
}

type entry struct {
	key    cache.Key
	value  cache.Value
	parent *list.Element
}

//
// Freq List: HEAD <-> Node <-> Node <-> Node <-> Tail
//                     Entry

func New(capacity uint64, maxElements uint64) *Cache {
	cache := &Cache{
		size:        0,
		capacity:    capacity,
		maxElements: maxElements,
		freq:        list.New(),
		caches:      make(map[cache.Key]*entry),
	}

	return cache
}

func (c *Cache) Set(key cache.Key, value cache.Value) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if e, ok := c.caches[key]; ok {
		e.value = value
		c.increment(e)
	} else {
		e := &entry{
			key:   key,
			value: value,
		}

		c.caches[key] = e
		c.increment(e)
	}

	c.checkCapacity()

}

func (c *Cache) Get(key cache.Key) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if e, ok := c.caches[key]; ok {
		c.increment(e)
		return e.value, true
	}

	return nil, false
}

func (c *Cache) Delete(key cache.Key) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if e, ok := c.caches[key]; ok {
		c.removeElement(e)
	}

}

func (c *Cache) Size() uint64 {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.size
}

func (c *Cache) ElementsCount() uint64 {
	c.mu.Lock()
	defer c.mu.Unlock()

	return uint64(len(c.caches))
}

func (c *Cache) Capacity() uint64 {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.capacity
}

func (c *Cache) SetCapacity(capacity uint64) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.capacity = capacity

	c.checkCapacity()
}

func (c *Cache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.freq.Init()
	c.size = 0
	c.caches = make(map[cache.Key]*entry)

}

func (c *Cache) Evict(n int) {
	if n <= 0 {
		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.evictElement(n)

}

func (c *Cache) removeElement(e *entry) {

	c.size -= uint64(c.caches[e.key].value.Size())

	n := e.parent
	items := n.Value.(*node).items

	delete(items, e.key)
	delete(c.caches, e.key)

	if len(items) == 0 {
		c.freq.Remove(n)
	}

}

func (c *Cache) evictElement(n int) {
	for ; n > 0 && len(c.caches) > 0; n-- {
		nn := c.freq.Front()

		for key, _ := range nn.Value.(*node).items {
			c.removeElement(c.caches[key])
			break
		}

	}

}

func (c *Cache) checkCapacity() {
	if c.capacity == 0 && c.maxElements == 0 {
		return
	}

	for (c.capacity != 0 && c.size > c.capacity) ||
		(c.maxElements != 0 && uint64(len(c.caches)) > c.maxElements) {
		c.evictElement(1)
	}
}

func (c *Cache) increment(e *entry) {
	var (
		freq    int
		n       *list.Element
		current *list.Element
	)

	current = e.parent

	if current == nil {
		freq = 1
		n = c.freq.Front()
	} else {
		freq = e.parent.Value.(*node).freq + 1
		n = current.Next()
	}

	if n == nil || n.Value.(*node).freq != freq {
		nn := &node{
			freq:  freq,
			items: make(map[cache.Key]uint8),
		}

		if current != nil {
			n = c.freq.InsertAfter(nn, current)
			delete(current.Value.(*node).items, e.key)
		} else {
			n = c.freq.PushFront(nn)
		}
	}

	e.parent = n

	n.Value.(*node).items[e.key] = 1

}
