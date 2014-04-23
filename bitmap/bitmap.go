package bitmap

import (
	"container/list"
	"sync"
)

type Bitmap struct {
	mu sync.Mutex

	size   int        // total elements stored
	pages  *list.List // @Todo: refactor this one
	option Option
}

type Option struct {
	AutoExpand  bool // Whether to alloc more spaces for coming elements
	AutoRecycle bool // Whether to recycle resources after delete elements from the bitmap
}

const (
	pageSize    = 1 << 12
	bitsPerByte = 1 << 3
	bitsPerPage = bitsPerByte * pageSize
)

type page struct {
	id   int
	bits []uint8
}

func NewOption() *Option {
	return &Option{
		AutoExpand:  false,
		AutoRecycle: false,
	}
}

func New(size int, option Option) *Bitmap {
	// For option, pass by value nor by pointer
	// Don't want the caller can modify the option after the init progress
	b := &Bitmap{
		option: option,
		pages:  list.New(),
		size:   0,
	}

	return b
}

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

func (b *Bitmap) getPage(n int, create bool) *page {
	pageIdx := b.getPageIndex(n)

	for e := b.pages.Front(); e != nil; {
		if e.Value.(*page).id == pageIdx {
			return e.Value.(*page)
		}

		if e.Next() != nil && e.Next().Value.(*page).id > pageIdx {
			if create {
				page := b.newPage(pageIdx)
				b.pages.InsertAfter(page, e)
				return page
			} else {
				return nil
			}
		}

		e = e.Next()
	}

	if create {
		page := b.newPage(pageIdx)
		b.pages.PushBack(page)
		return page
	} else {
		return nil
	}
}

func (b *Bitmap) setBitInPage(n int, set bool) {
	page := b.getPage(n, true)
	idx := (n - 1) % bitsPerPage / bitsPerByte
	if set {
		page.bits[idx] |= 1 << uint8((n-1)%bitsPerPage%bitsPerByte)
	} else {
		page.bits[idx] &^= 1 << uint8((n-1)%bitsPerPage%bitsPerByte)
	}
}

func (b *Bitmap) Check(n int) bool {

	page := b.getPage(n, false)
	if page == nil {
		return false
	}

	idx := (n - 1) % bitsPerPage / bitsPerByte
	return page.bits[idx]&(1<<uint8((n-1)%bitsPerPage%bitsPerByte)) > 0
}

func (b *Bitmap) Clear(n int) {
	b.setBitInPage(n, false)
	b.size--
}

func (b *Bitmap) Set(n int) {
	b.setBitInPage(n, true)
	b.size++
}

func (b *Bitmap) Size() int {
	return b.size
}

func (b *Bitmap) Gc() {
}

func (b *Bitmap) Capacity() int {
}
