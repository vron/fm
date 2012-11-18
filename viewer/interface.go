package main

// TODO: Handle async in this interface
import (
	"fmt"
)

type node [3]float64
type Point [3]float64
type element interface {
	Nodes() []uint32
}

type quad struct {
	nodes [4]uint32
}

func (q *quad) Nodes() []uint32 {
	return q.nodes[:]
}

type tria struct {
	nodes [3]uint32
}

func (t *tria) Nodes() []uint32 {
	return t.nodes[:]
}

var (
	nodes           = map[uint32]node{}
	itemid   uint32 = 0
	elements        = map[uint32]element{}
)

// Create a node at the provided point and return the node identifier, else 0 and a error
// It is guareneed that the index of the new node is the last index + 1, and that the first node
// always is 1, this way a file may use absolute indexes without storing the node number in a variable
func Node(x, y, z float64) (uint32, error) {
	itemid++
	nodes[itemid] = node{x, y, z}
	return itemid, nil
}

// Create a quadrilatilar element and return id
func Quad(n1, n2, n3, n4 uint32) (uint32, error) {
	q := quad{[4]uint32{n1, n2, n3, n4}}
	for i, v := range q.nodes {
		if _, ok := nodes[v]; !ok {
			return 0, fmt.Errorf("Argument #%d (id %d) not a node", i, v)
		}
	}
	itemid++
	elements[itemid] = &q
	return itemid, nil
}

// Create triangulat element and return id
func Tria(n1, n2, n3 uint32) (uint32, error) {
	t := tria{[3]uint32{n1, n2, n3}}
	for i, v := range t.nodes {
		if _, ok := nodes[v]; !ok {
			return 0, fmt.Errorf("Argument #%d (id %d) not a node", i, v)
		}
	}
	itemid++
	elements[itemid] = &t
	return itemid, nil
}
