package roundrobin

import (
	"github.com/flatpeach/coconut/util"
)

type RoundRobin struct {
	nodes []Node

	index  int
	weight int
}

type Node interface {
	Weight() int
}

func New(nodes ...Node) *RoundRobin {
	rr := &RoundRobin{
		index:  -1,
		weight: 0,
		nodes:  make([]Node, len(nodes)),
	}

	for idx, node := range nodes {
		rr.nodes[idx] = node
	}

	return rr
}

func (rr *RoundRobin) gcd() uint {
	weights := make([]uint, len(rr.nodes))

	for idx, node := range rr.nodes {
		weights[idx] = uint(node.(Node).Weight())
	}

	return util.Gcd(weights...)
}

func (rr *RoundRobin) maxWeight() int {
	weight := 0

	for _, node := range rr.nodes {
		if node.(Node).Weight() > weight {
			weight = node.(Node).Weight()
		}
	}

	return weight
}

func (rr *RoundRobin) Next() Node {
	for {
		rr.index = (rr.index + 1) % len(rr.nodes)

		if rr.index == 0 {
			rr.weight = rr.weight - int(rr.gcd())

			if rr.weight <= 0 {
				rr.weight = rr.maxWeight()

				if rr.weight == 0 {
					return nil
				}
			}
		}

		if (rr.nodes[rr.index]).Weight() >= rr.weight {
			return rr.nodes[rr.index]
		}

	}
}
