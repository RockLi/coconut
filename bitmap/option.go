// Copyright 2014 The coconut Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bitmap

// Option used to construct the Bitmap
type Option struct {
	// Automatically to alloc more spaces for coming elements
	AutoExpand bool

	// Automatically to recycle resources after delete elements from the bitmap
	AutoRecycle bool

	// Initial capacity of this Bitmap
	Capacity int
}

func (o *Option) clone() *Option {
	return &Option{
		AutoExpand:  o.AutoExpand,
		AutoRecycle: o.AutoRecycle,
		Capacity:    o.Capacity,
	}
}
