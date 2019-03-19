package llb2dot

import (
	"gonum.org/v1/gonum/graph/simple"
)

// LLB2Graph convert llb DAG graph to dot language.
func LLB2Graph(ops LLBOps) (*simple.DirectedGraph, error) {
	g := simple.NewDirectedGraph()
	nm := newNodeManager()

	for _, op := range ops {
		to := nm.createIfNotExists(g, string(op.Digest), op.getDesc())
		for _, input := range op.Op.Inputs {
			from := nm.createIfNotExists(g, string(input.Digest), op.getDesc())
			e := g.NewEdge(from, to)
			g.SetEdge(e)
		}
	}
	return g, nil
}
