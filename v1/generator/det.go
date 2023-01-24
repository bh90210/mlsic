// Package generator implements mlsic.Graph but instead of using
// GNNs for graph generation, it uses a deterministic algos.
// It exists mostly as a developer's tool to help seed the
// process of writing Algo1.
package generator

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"

	"github.com/bh90210/mlsic/v1"
	"github.com/bh90210/mlsic/v1/internal"
	"gonum.org/v1/gonum/graph/encoding"
	"gonum.org/v1/gonum/graph/encoding/dot"
	"gonum.org/v1/gonum/graph/simple"
)

var _ mlsic.Graph = (*Det)(nil)

// Det implements mlsic.VGAE.
type Det struct {
	// TotalGraphs total numbers of graphs to be generated.
	TotalGraphs int
	// Print if true renders all graphs generated as .dot and subsequently .svg files.
	Print bool
	// Seed of the rand.Seed() function.
	Seed int64
}

// Dump returns a deterministic sequence of graphs for dev purposes.
func (d *Det) Dump() ([]*simple.WeightedUndirectedGraph, error) {
	// Init graphs variable. This slice will hold and return all generated graphs.
	var graphs []*simple.WeightedUndirectedGraph

	// Seed the rand function.
	rand.Seed(d.Seed)

	// Start the creation of graphs.
	for i := 0; i < d.TotalGraphs; i++ {
		g := simple.NewWeightedUndirectedGraph(0, 0)
		totalNodes := rand.Intn(100)
		for {
			if totalNodes > 10 {
				break
			}

			totalNodes = rand.Intn(100)
		}

		// Create the actual nodes.
		for i := 0; i < totalNodes; i++ {
			g.AddNode(g.NewNode())
		}

		// Get nodes as slice.
		// nodes := graph.NodesOf(g.Nodes())

		// Second mirror connect nodes.
		mirror := totalNodes
		// If mirror (total nodes) is odd.
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
					Dot: []encoding.Attribute{{Key: "label", Value: strconv.Itoa(1)}},
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
						Dot: []encoding.Attribute{{Key: "label", Value: strconv.Itoa(2)}},
					})
				}
			}
		}

		graphs = append(graphs, g)
	}

	if d.Print {
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
