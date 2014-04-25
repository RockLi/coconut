// Copyright 2014 The coconut Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package cache provides interfaces for cache functions
package cache

// Cache is the common interface implemented by all cache functions
type Cache interface {
	Set(k Key, d Data)

	Get(k Key) (Data, bool)

	Delete(k Key)

	// Size returns the count in bytes of the cache
	Size() uint64

	ElementsCount() uint64

	Clear()

	Evict(n int)

	Capacity() uint64

	SetCapacity(cap uint64)

	// Full returns whether current cache is full or not
	Full() bool
}

type Key interface{}

type Option interface {

	// Capacity returns the current max allowed bytes
	// to store in the cache
	// Zero means no limitation.
	Capacity() uint64

	// SetCapacity sets the max allowed bytes
	// to store in the cache
	// Zero means no limitation.
	SetCapacity(cap uint64)

	// MaxElements returns current max elements totally
	// to store in the cache
	// Zero means no limitation.
	MaxElements() uint64

	// SetMaxElements sets the max count of elements
	// to store in the cache.
	// Zero means no limitation.
	SetMaxElements(count uint64)
}

// Data is the interface to insert data to the cache
type Data interface {

	// Size returns the count in bytes of this data
	Size() uint64
}
