package cache

type Key string

type Value interface {
	Size() int
}

type Cache interface {
	Set(key Key, value Value)

	Get(key Key) (interface{}, bool)

	Delete(key Key)

	Size() uint64

	ElementsCount() uint64

	Capacity() uint64

	SetCapacity() uint64

	Clear()

	Evict(n int)
}
