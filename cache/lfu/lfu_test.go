// Copyright 2014 The coconut Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lfu

import (
	"testing"
)

type cacheItem struct {
	v []byte
}

func (i *cacheItem) Size() uint64 {
	return uint64(len(i.v))
}

func TestLFU(t *testing.T) {
	c := New(&Option{0, 2})

	key1 := "key1"
	key2 := "key2"
	key3 := "key3"

	value1 := &cacheItem{[]byte("value1")}
	value2 := &cacheItem{[]byte("value2")}
	value3 := &cacheItem{[]byte("value3")}

	c.Set(key1, value1)
	c.Set(key2, value2)

	c.Get(key1)
	c.Get(key2)

	c.Set(key3, value3)

	if c.ElementsCount() != 2 {
		t.Fatal("Total elements should equal to 2")
	}

	v, ok := c.Get(key1)
	if !ok {
		t.Fatal("Key1 should in the cache")
	}

	if v.(*cacheItem) != value1 {
		t.Fatal("Value1 changed in cache")
	}

	v, ok = c.Get(key2)
	if !ok {
		t.Fatal("Key2 should in the cache")
	}

	if v.(*cacheItem) != value2 {
		t.Fatal("Value2 changed in cache")
	}

	c.Get(key1)

	c.Evict(1)

	if c.ElementsCount() != 1 {
		t.Fatal("Current should have only one element ")
	}

	v, ok = c.Get(key1)
	if !ok {
		t.Fatal("Now should only key1 in the cache")
	}

	c.Clear()

	if c.ElementsCount() != 0 {
		t.Fatal("Should be a empty cache")
	}

}
