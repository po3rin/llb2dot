package llb2dot

import (
	"io"

	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/encoding/dot"
	"gonum.org/v1/gonum/graph/simple"
)

type node struct {
	id   int64
	desc string
}

func (n node) ID() int64 {
	return n.id
}

func (n node) DOTID() string {
	return n.desc
}

type nodeManager struct {
	increment  int64
	addednodes map[string]int64
}

func newNodeManager() *nodeManager {
	return &nodeManager{
		addednodes: map[string]int64{},
	}
}

func (nm *nodeManager) createIfNotExists(g *simple.DirectedGraph, digest, description string) graph.Node {
	var n graph.Node
	id, ok := nm.addednodes[digest]
	if !ok {
		n = nm.create(g, digest, description)
		g.AddNode(n)
	} else {
		n = g.Node(id)
	}
	return n
}

func (nm *nodeManager) create(g *simple.DirectedGraph, nodeDigest, desc string) node {
	nm.increment++
	nm.addednodes[nodeDigest] = nm.increment
	return node{id: nm.increment, desc: desc}
}

// WriteDOT output graph to dot language.
func WriteDOT(w io.Writer, g graph.Graph) error {
	b, err := dot.Marshal(g, "llb", "", "")
	if err != nil {
		return err
	}
	_, err = w.Write(b)
	return err
}
