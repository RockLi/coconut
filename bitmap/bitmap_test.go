// Copyright 2014 The coconut Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bitmap

import (
	"testing"
)

func TestBitmap(t *testing.T) {
	bm := New(nil)

	bm.Set(10)

	if bm.Size() != 1 {
		t.Fatal("Bitmap should has one element")
	}

	if bm.Capacity() != bitsPerPage {
		t.Fatal("Capacity should be 4K * 8(default)")
	}

	bm.Clear(10)
	if bm.Size() != 0 {
		t.Fatal("failed to clean")
	}

	bm.Set(10)
	bm.ClearAll()
	if bm.Capacity() != bitsPerPage {
		t.Fatal("Clear action should cleared the whole memory")
	}

	for i := 1; i < 150000; i++ {
		bm.Set(i)

		if !bm.Test(i) {
			t.Fatalf("%d should has one element", i)
		}

		if bm.Test(i + 1) {
			t.Fatalf("%d should not exist", i+1)
		}

	}

	for i := 1; i < 150000; i++ {
		bm.Clear(i)

		if bm.Test(i) {
			t.Fatalf("%d should not exist", i)
		}

	}

	bm.Gc()

	if bm.Capacity() != 0 {
		t.Fatal("No resident page should be found")
	}

	if bm.Size() != 0 {
		t.Fatal("Should be empty BitMap")
	}

}

func TestBitMapAutoExpandDisabled(t *testing.T) {
	o := NewOption(10, false, false)

	bm := New(o)

	bm.Set(100) // Should fail at this step
	if bm.Test(100) {
		t.Fatal("autoexpand should be disabled")
	}

	bm.Set(4)
	if !bm.Test(4) {
		t.Fatal("4 should in this bitmap")
	}
}

func TestBitMapAutoExpandEnabled(t *testing.T) {
	o := NewOption(10, true, false)
	bm := New(o)

	bm.Set(100000)
	if !bm.Test(100000) {
		t.Fatal("Autoexpand doesn't work")
	}

}

func TestBitMapAutoRecycleEnabled(t *testing.T) {
	o := NewOption(10, true, true)
	bm := New(o)

	bm.Set(100000)
	bm.Clear(100000)

	if bm.pages.Len() != 0 {
		t.Fatal("Autorecycle doesn't work")
	}

}

func TestBitMapAutoRecycleDisabled(t *testing.T) {
	o := NewOption(10, true, false)
	bm := New(o)

	bm.Set(100000)
	bm.Clear(100000)

	if bm.pages.Len() != 1 {
		t.Fatal("Autorecycle doesn't work")
	}

	bm.Gc()

	if bm.pages.Len() != 0 {
		t.Fatal("Failed to run GC")
	}

}
