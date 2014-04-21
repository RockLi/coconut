package lru

import (
	"github.com/flatpeach/coconut/cache"
	"testing"
)

type CacheItem struct {
	v []byte
}

func (i *CacheItem) Size() int {
	return len(i.v)
}

func TestLRUBasic(t *testing.T) {
	c := New(1<<20, 0)

	if c.Capacity() != 1<<20 {
		t.Fatal("The capacity of LRU cache not matched!")
	}

	key := cache.Key("hello")

	value := &CacheItem{[]byte("HelloWorld")}

	c.Set(key, value)

	v, ok := c.Get(key)

	if !ok {
		t.Fatal("Didn't get the key in memory")
		return
	}

	if v.(*CacheItem) != value {
		t.Fatal("Data mismatched!")
		return
	}

	if c.Size() != uint64(value.Size()) {
		t.Fatal("size not matched")
	}

	c.SetCapacity(1 << 30)

	if c.Capacity() != 1<<30 {
		t.Fatal("Set capacity failed")
	}

	c.Delete(key)

	if c.Size() != 0 {
		t.Fatal("Failed to delete one item")
	}

	v, ok = c.Get(key)
	if ok {
		t.Fatal("Failed to delete one element, still can access")
	}

}

func TestLRUEvict(t *testing.T) {
	c := New(0, 2)

	key1 := cache.Key("k1")
	key2 := cache.Key("k2")
	key3 := cache.Key("k3")

	val1 := &CacheItem{[]byte("HelloWorld")}
	val2 := &CacheItem{[]byte("HelloWorld")}
	val3 := &CacheItem{[]byte("HelloWorld")}

	c.Set(key1, val1)
	if c.ElementsCount() != 1 {
		t.Fatal("Count of elements should be equal to 1")
	}

	c.Set(key2, val2)
	if c.ElementsCount() != 2 {
		t.Fatal("Count of elements should be equal to 2")
	}

	c.Set(key3, val3)
	if c.ElementsCount() != 2 {
		t.Fatal("Count of elements should be equal to 2")
	}

	c.Evict(1)
	if c.ElementsCount() != 1 {
		t.Fatal("Count of elements should be equal to 1")
	}

	c.Evict(1)
	if c.ElementsCount() != 0 {
		t.Fatal("Count of elements should be equal to 0")
	}

}
