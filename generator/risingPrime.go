// Package generator implements mlsic.Graph but instead of using
// GNNs for graph generation, it uses deterministic algos.
// It exists mostly as a developer's tool to help seed the
// process of writing Algo1.
package generator

import (
	"fmt"
	"math/big"
	"os"
	"os/exec"
	"strconv"

	"github.com/bh90210/mlsic"
	"github.com/bh90210/mlsic/internal"
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/encoding"
	"gonum.org/v1/gonum/graph/encoding/dot"
	"gonum.org/v1/gonum/graph/simple"
)

var _ mlsic.Graph = (*Prime1)(nil)

// Prime1 implements mlsic.Graph. If it a deterministic graph generator
// that starts from a asimple graph with few nodes and scales up nodes
// for the next graph following primes numbers.
type Prime1 struct {
	// TotalGraphs total numbers of graphs to be generated.
	TotalGraphs int
	// Print if true renders all graphs generated as .dot and subsequently .svg files.
	Print bool
}

// Dump returns a deterministic sequence of graphs for dev purposes.
func (p *Prime1) Dump() ([]*simple.WeightedUndirectedGraph, error) {
	// Init graphs variable. This slice will hold and return all generated graphs.
	var graphs []*simple.WeightedUndirectedGraph

	// Primes is a helper slice with equal primes as p.MaxNodes.
	var primes []int
	for i := 13; ; i++ {
		if len(primes) == p.TotalGraphs {
			break
		}

		if big.NewInt(int64(i)).ProbablyPrime(i) {
			primes = append(primes, i)
		}
	}

	// Start the creation of graphs.
	for i := 0; i < p.TotalGraphs; i++ {
		g := simple.NewWeightedUndirectedGraph(0, 0)
		totalNodes := primes[i]

		// Create the actual nodes.
		for i := 0; i < totalNodes; i++ {
			g.AddNode(g.NewNode())
		}

		// Get nodes as slice.
		nodes := graph.NodesOf(g.Nodes())

		// First connect prime nodes with each other with weight 1.
		for _, node := range nodes {
			if big.NewInt(node.ID()).ProbablyPrime(int(node.ID())) {
				for _, n := range nodes {
					// Skip self.
					if node.ID() == n.ID() {
						continue
					}

					// If other node is a prime try connecting them.
					if big.NewInt(n.ID()).ProbablyPrime(int(n.ID())) {
						edge := g.Edge(node.ID(), n.ID())
						if edge == nil {
							g.SetWeightedEdge(&internal.AttrEdge{
								F:   node,
								T:   n,
								W:   float64(1),
								Dot: []encoding.Attribute{{Key: "label", Value: strconv.Itoa(1)}},
							})
						}
					}
				}
			}
		}

		// Second mirror connect nodes.
		mirror := totalNodes
		// If mirror (total nodes) is odp.
		if mirror%2 == 1 {
			// Make it even.
			mirror--
		}

		half := mirror / 2

		for i := 1; i < half; i++ {
			edge := g.Edge(int64(i), int64(i+half))
			if edge == nil {
				g.SetWeightedEdge(&internal.AttrEdge{
					F:   g.Node(int64(i)),
					T:   g.Node(int64(i + half)),
					W:   float64(2),
					Dot: []encoding.Attribute{{Key: "label", Value: strconv.Itoa(2)}},
				})
			}
		}

		// Third connect next node.
		for i := 0; i < totalNodes; i++ {
			if i+1 < totalNodes {
				edge := g.Edge(int64(i), int64(i+1))
				if edge == nil {
					g.SetWeightedEdge(&internal.AttrEdge{
						F:   g.Node(int64(i)),
						T:   g.Node(int64(i + 1)),
						W:   float64(3),
						Dot: []encoding.Attribute{{Key: "label", Value: strconv.Itoa(3)}},
					})
				}
			}
		}

		graphs = append(graphs, g)
	}

	if p.Print {
		for i, g := range graphs {
			// Create a graphviz dot file.
			data, err := dot.Marshal(g, "", "", "  ")
			if err != nil {
				return nil, err
			}

			err = os.WriteFile(fmt.Sprint(i)+".dot", data, 0644)
			if err != nil {
				return nil, err
			}

			_, err = exec.Command("dot", "-Tsvg", fmt.Sprint(i)+".dot", "-o", fmt.Sprint(i)+".svg").Output()
			// dot -Tsvg '/current/path/test.dot' -o output.svg
			if err != nil {
				return nil, err
			}
		}
	}

	return graphs, nil
}
