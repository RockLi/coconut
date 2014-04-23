package lru

// LRU package implements a lightweight LRU cache, it's threadsafe

import (
	"container/list"
	"sync"
)

type Value interface {
	Size() int
}

type entry struct {
	key   string // Key for this item
	value Value  // Value for this item
}

type Cache struct {
	mu sync.Mutex

	size        uint64
	capacity    uint64
	maxElements uint64
	items       *list.List
	caches      map[string]*list.Element
}

// Both capacity and maxElements equal to 0 means no limitation
func New(capacity uint64, maxElements uint64) *Cache {
	c := &Cache{
		size:        0,
		capacity:    capacity,
		maxElements: maxElements,
		items:       list.New(),
		caches:      make(map[string]*list.Element),
	}

	return c
}

func (c *Cache) Set(key string, value Value) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if e, ok := c.caches[key]; ok {
		c.items.MoveToFront(e)
		e.Value.(*entry).value = value
		return
	}

	item := &entry{
		key:   key,
		value: value,
	}

	e := c.items.PushFront(item)
	c.caches[key] = e

	c.size += uint64(value.Size())

	c.checkCapacity()
}

func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if e, ok := c.caches[key]; ok {
		c.items.MoveToFront(e)
		return e.Value.(*entry).value, true
	}

	return nil, false
}

func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if e, ok := c.caches[key]; ok {
		c.size -= uint64(e.Value.(*entry).value.Size())
		c.items.Remove(e)
		delete(c.caches, key)
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

	c.items.Init()
	c.caches = make(map[string]*list.Element)
	c.size = 0
}

func (c *Cache) checkCapacity() {
	// Do not lock here

	if c.capacity == 0 && c.maxElements == 0 {
		return
	}

	for (c.capacity != 0 && c.size > c.capacity) ||
		(c.maxElements != 0 && uint64(len(c.caches)) > c.maxElements) {
		c.evictElement(1)
	}
}

func (c *Cache) evictElement(n int) {
	// Do not lock here
	for ; n > 0 && len(c.caches) > 0; n-- {

		e := c.items.Back()
		v := e.Value.(*entry)

		c.items.Remove(e)
		delete(c.caches, v.key)

		c.size -= uint64(v.value.Size())

	}

}

// Evict n elements from LRU cache
func (c *Cache) Evict(n int) {
	if n <= 0 {
		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.evictElement(n)

}
