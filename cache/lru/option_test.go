// Copyright 2014 The coconut Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lru

import (
	"testing"
)

func TestOption(t *testing.T) {

	o := new(Option)

	if o.Capacity() != 0 {
		t.Fatal("Default init of option is not correct")
	}

	if o.MaxElements() != 0 {
		t.Fatal("Default init of option is not correct")
	}

	o.SetCapacity(100)
	if o.Capacity() != 100 {
		t.Fatal("Failed to adjust the capacity")
	}

	o.SetMaxElements(1000)
	if o.MaxElements() != 1000 {
		t.Fatal("Failed to adjust the max elements")
	}

}
