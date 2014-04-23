package bitmap

import (
	_ "fmt"
	"testing"
)

func TestBitmap(t *testing.T) {

	bm := New(100, *NewOption())

	// bm.Set(10)

	// if bm.Size() != 1 {
	// 	t.Fatal("Bitmap should has one element")
	// }

	// bm.Clear(10)

	if bm.Size() != 0 {
		t.Fatal("failed to clean")
	}

	for i := 1; i < 150000; i++ {
		bm.Set(i)

		if !bm.Check(i) {
			t.Fatalf("%d should has one element", i)
		}

		if bm.Check(i + 1) {
			t.Fatalf("%d should not exist", i+1)
		}

	}

}
