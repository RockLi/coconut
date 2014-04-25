// Copyright 2014 The coconut Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package lru provides one implementation of the LRU cache
package lru

import (
	"container/list"
	"github.com/flatpeach/coconut/cache"
	"sync"
)

type entry struct {
	key  cache.Key  // Key for this item
	data cache.Data // Data for this item
}

type Cache struct {
	mu sync.Mutex

	size   uint64
	items  *list.List
	caches map[cache.Key]*list.Element
	o      *Option
}

func New(o *Option) *Cache {
	c := &Cache{
		size:   0,
		items:  list.New(),
		caches: make(map[cache.Key]*list.Element),
	}

	if o == nil {
		c.o = &Option{0, 0}
	} else {
		c.o = &Option{o.Capacity(), o.MaxElements()} // copy by value
	}

	return c
}

func (c *Cache) Set(key cache.Key, data cache.Data) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if e, ok := c.caches[key]; ok {
		c.items.MoveToFront(e)
		e.Value.(*entry).data = data
		return
	}

	item := &entry{
		key:  key,
		data: data,
	}

	e := c.items.PushFront(item)
	c.caches[key] = e

	c.size += uint64(data.Size())
	c.checkCapacity()
}

func (c *Cache) Get(key cache.Key) (cache.Data, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if e, ok := c.caches[key]; ok {
		c.items.MoveToFront(e)
		return e.Value.(*entry).data, true
	}

	return nil, false
}

func (c *Cache) Delete(key cache.Key) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if e, ok := c.caches[key]; ok {
		c.size -= uint64(e.Value.(*entry).data.Size())
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

	return c.o.capacity
}

func (c *Cache) SetCapacity(capacity uint64) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.o.capacity = capacity
	c.checkCapacity()
}

func (c *Cache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items.Init()
	c.caches = make(map[cache.Key]*list.Element)
	c.size = 0
}

func (c *Cache) checkCapacity() {
	if c.o.capacity == 0 && c.o.maxElements == 0 {
		return
	}

	for (c.o.capacity != 0 && c.size > c.o.capacity) ||
		(c.o.maxElements != 0 && uint64(len(c.caches)) > c.o.maxElements) {
		c.evictElement(1)
	}
}

func (c *Cache) evictElement(n int) {
	for ; n > 0 && len(c.caches) > 0; n-- {
		e := c.items.Back()
		v := e.Value.(*entry)

		c.items.Remove(e)
		delete(c.caches, v.key)

		c.size -= uint64(v.data.Size())
	}
}

func (c *Cache) Evict(n int) {
	if n <= 0 {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()

	c.evictElement(n)

}

func (c *Cache) Full() bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.o.capacity == 0 && c.o.maxElements == 0 {
		return false
	}

	if (c.o.capacity != 0 && c.size >= c.o.capacity) ||
		(c.o.maxElements != 0 && uint64(len(c.caches)) >= c.o.maxElements) {
		return true
	}

	return false

}
