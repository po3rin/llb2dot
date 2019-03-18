package llb2dot

import (
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/simple"
)

func newNodeIfNotExits(g *simple.DirectedGraph, digest, description string) graph.Node {
	var n graph.Node
	id, ok := addednodes[digest]
	if !ok {
		n = newNode(g, digest, description)
		g.AddNode(n)
	} else {
		n = g.Node(id)
	}
	return n
}

// LLB2Graph convert llb DAG graph to dot language.
func LLB2Graph(llbOps LLBOps) (*simple.DirectedGraph, error) {
	g := simple.NewDirectedGraph()
	for _, llbOp := range llbOps {
		to := newNodeIfNotExits(g, string(llbOp.Digest), llbOp.getDesc())
		for _, input := range llbOp.Op.Inputs {
			from := newNodeIfNotExits(g, string(input.Digest), llbOp.getDesc())
			e := g.NewEdge(from, to)
			g.SetEdge(e)
		}
	}
	return g, nil
}
