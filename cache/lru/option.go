// Copyright 2014 The coconut Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lru

type Option struct {
	capacity    uint64
	maxElements uint64
}

func (o *Option) Capacity() uint64 {
	return o.capacity
}

func (o *Option) SetCapacity(cap uint64) {
	o.capacity = cap
}

func (o *Option) MaxElements() uint64 {
	return o.maxElements
}

func (o *Option) SetMaxElements(count uint64) {
	o.maxElements = count
}
