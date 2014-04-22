package roundrobin

import (
	"strconv"
	"strings"
	"testing"
)

type MyServer struct {
	id     int
	weight int
}

func (s MyServer) Weight() int {
	return s.weight
}

func TestRR(t *testing.T) {
	rr := New(
		MyServer{1, 1},
		MyServer{2, 1},
		MyServer{3, 1},
		MyServer{4, 1},
	)

	seq := ""
	for i := 0; i < 8; i++ {
		node := rr.Next()
		seq += strconv.Itoa(node.(MyServer).id)
	}

	if seq != "12341234" {
		t.Fatal("Sequence should be 12341234")
	}

}

func TestWeightedRR(t *testing.T) {
	rr := New(
		MyServer{1, 10},
		MyServer{2, 50},
		MyServer{3, 50},
		MyServer{4, 10},
	)

	seq := ""
	for i := 0; i < 100; i++ {
		node := rr.Next()
		seq += strconv.Itoa(node.(MyServer).id)
	}

	if !strings.ContainsAny(seq, "1") ||
		!strings.Contains(seq, "2") ||
		!strings.Contains(seq, "3") ||
		!strings.Contains(seq, "4") {
		t.Fatal("every node should has the opportunity ")
	}

}
