// Copyright 2014 The coconut Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
	Package bitmap provides a sparse bitmap implementation.
*/
package bitmap

import (
	"container/list"
	"os"
	"sync"
)

type Bitmap struct {
	mu sync.Mutex

	size   int        // total elements stored
	pages  *list.List // @Todo: here can be optimized
	option *Option
}

// Option used to construct the Bitmap
type Option struct {
	// Automatically to alloc more spaces for coming elements
	AutoExpand bool

	// Automatically to recycle resources after delete elements from the bitmap
	AutoRecycle bool

	// Initial capacity of this Bitmap
	Capacity int
}

// return a new option
func (o *Option) clone() *Option {
	return &Option{
		AutoExpand:  o.AutoExpand,
		AutoRecycle: o.AutoRecycle,
		Capacity:    o.Capacity,
	}
}

const (
	bitsPerByte = 1 << 3 // default to 8bits per byte, only consider modern computing
)

var (
	pageSize    = os.Getpagesize()
	bitsPerPage = bitsPerByte * pageSize
)

// one page is the minimal unit for managing lot's of bits
type page struct {
	id   int
	size int
	bits []uint8
}

func NewOption(capacity int, autoExpand, autoRecycle bool) *Option {
	return &Option{
		AutoExpand:  autoExpand,
		AutoRecycle: autoRecycle,
		Capacity:    capacity,
	}
}

// Return back a new Bitmap according the option passed in
func New(option *Option) *Bitmap {
	if option == nil {
		option = NewOption(0, true, true)
	}
	b := &Bitmap{
		option: option.clone(),
		pages:  list.New(),
		size:   0,
	}

	return b
}

// alloc one new page
func (b *Bitmap) newPage(id int) *page {
	p := &page{
		id:   id,
		bits: make([]byte, pageSize),
	}

	return p
}

func (b *Bitmap) getPageIndex(n int) int {
	if n%bitsPerPage != 0 {
		return n / bitsPerPage
	} else {
		return n/bitsPerPage - 1
	}
}

func (b *Bitmap) getPage(n int, create bool) *list.Element {
	if n > b.option.Capacity && !b.option.AutoExpand {
		return nil
	}

	pageIdx := b.getPageIndex(n)

	for e := b.pages.Front(); e != nil; {
		if e.Value.(*page).id == pageIdx {
			return e
		}

		if e.Next() != nil && e.Next().Value.(*page).id > pageIdx {
			if !create {
				return nil
			}

			page := b.newPage(pageIdx)
			b.pages.InsertAfter(page, e)
			return e.Next()

		}

		e = e.Next()
	}

	if !create {
		return nil
	}

	page := b.newPage(pageIdx)
	b.pages.PushBack(page)

	return b.pages.Back()
}

func (b *Bitmap) setBitInPage(n int, set bool) {
	e := b.getPage(n, true)
	if e == nil {
		// silently ignored this request
		return
	}

	page := e.Value.(*page)

	idx := (n - 1) % bitsPerPage / bitsPerByte
	if set {
		page.bits[idx] |= 1 << uint8((n-1)%bitsPerPage%bitsPerByte)
		page.size++
	} else {
		page.bits[idx] &^= 1 << uint8((n-1)%bitsPerPage%bitsPerByte)
		page.size--

		if page.size == 0 && b.option.AutoRecycle {
			b.pages.Remove(e)
		}

	}
}

// Test whether one bit is set or not
func (b *Bitmap) Test(n int) bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	e := b.getPage(n, false)
	if e == nil {
		return false
	}

	page := e.Value.(*page)

	idx := (n - 1) % bitsPerPage / bitsPerByte
	return page.bits[idx]&(1<<uint8((n-1)%bitsPerPage%bitsPerByte)) > 0
}

// Clear one bit
func (b *Bitmap) Clear(n int) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.setBitInPage(n, false)
	b.size--
}

// Set one bit
func (b *Bitmap) Set(n int) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.setBitInPage(n, true)
	b.size++
}

// Reinit the whole Bitmap
func (b *Bitmap) ClearAll() {
	b.mu.Lock()
	defer b.mu.Unlock()

	var p *page

	for e := b.pages.Front(); e != nil; e = e.Next() {
		p = e.Value.(*page)

		// Actually I just want memset, painful, waiting for new Go release
		for i := range p.bits {
			p.bits[i] = 0
		}

		p.size = 0

	}

	b.size = 0

}

// Total count of bits setted
func (b *Bitmap) Size() int {
	b.mu.Lock()
	defer b.mu.Unlock()

	return b.size
}

// Recycle unused pages
func (b *Bitmap) Gc() {
	b.mu.Lock()
	defer b.mu.Unlock()

	for e := b.pages.Front(); e != nil; {

		if e.Value.(*page).size == 0 {
			old := e.Next()
			b.pages.Remove(e)
			e = old
		} else {
			e = e.Next()
		}

	}

}

// How many bits can be set in this bitmap
func (b *Bitmap) Capacity() int {
	b.mu.Lock()
	defer b.mu.Unlock()

	if !b.option.AutoExpand {
		return b.option.Capacity
	}

	e := b.pages.Back()
	if e == nil {
		return 0
	}

	return (e.Value.(*page).id + 1) * bitsPerPage
}
