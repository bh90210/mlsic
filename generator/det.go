// Package generator implements mlsic.Graph but instead of using
// GNNs for graph generation, it uses deterministic algos.
// It exists mostly as a developer's tool to help seed the
// process of writing Algo1.
package generator

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strconv"

	"github.com/bh90210/mlsic"
	"github.com/bh90210/mlsic/internal"
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/encoding"
	"gonum.org/v1/gonum/graph/encoding/dot"
	"gonum.org/v1/gonum/graph/network"
	"gonum.org/v1/gonum/graph/path"
	"gonum.org/v1/gonum/graph/simple"
)

var _ mlsic.Graph = (*Det)(nil)

// Det implements mlsic.Graph. It is a deterministic graph generator
// meant to be used for developing Algo1.
type Det struct {
	// TotalGraphs total numbers of graphs to be generated.
	TotalGraphs int
	// Print if true renders all graphs generated as .dot and subsequently .svg files.
	Print bool
	// Seed of the rand.Seed() function.
	Seed int64
	// MaxNodes is the maximum number of nodes allowed per graph.
	MaxNodes int
	// MaxEdges is the maximum number of edges a node can have.
	MaxEdges int
	// MaxWeight is the maximum number an edge can be assigned.
	MaxWeight int
}

// Dump returns a deterministic sequence of graphs for dev purposes.
func (d *Det) Dump() ([]*simple.WeightedUndirectedGraph, error) {
	// Init graphs variable. This slice will hold and return all generated graphs.
	var graphs []*simple.WeightedUndirectedGraph

	// Seed the rand function.
	rand.Seed(d.Seed)

	// Start the creation of graphs.
	for y := 0; y < d.TotalGraphs; y++ {
		g := simple.NewWeightedUndirectedGraph(0, 0)

		// Calculate total nodes number.
		totalNodes := rand.Intn(d.MaxNodes)
		for {
			if totalNodes > 10 {
				break
			}

			totalNodes = rand.Intn(d.MaxNodes)
		}

		switch y {
		case 0:
			// Create the actual nodes.
			for i := 0; i < totalNodes; i++ {
				g.AddNode(g.NewNode())
			}

			// Get nodes as slice.
			// nodes := graph.NodesOf(g.Nodes())
			// for _, node := range nodes {
			// 	if big.NewInt(node.ID()).ProbablyPrime(int(node.ID())) {
			// 		for _, n := range nodes {
			// 			// Skip self.
			// 			if node.ID() == n.ID() {
			// 				continue
			// 			}

			// 			// If other node is a prime try connecting them.
			// 			if big.NewInt(n.ID()).ProbablyPrime(int(n.ID())) {
			// 				edge := g.Edge(node.ID(), n.ID())
			// 				if edge == nil {
			// 					g.SetWeightedEdge(&internal.AttrEdge{
			// 						F:   node,
			// 						T:   n,
			// 						W:   float64(1),
			// 						Dot: []encoding.Attribute{{Key: "label", Value: strconv.Itoa(1)}},
			// 					})
			// 				}
			// 			}
			// 		}
			// 	}
			// }

			// Second mirror connect nodes.
			mirror := totalNodes
			// If mirror (total nodes) is odd.
			if mirror%2 == 1 {
				// Make it even.
				mirror--
			}

			half := mirror / 2

			// We ignore node 0.
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

		default:
			graph.CopyWeighted(g, graphs[y-1])

			bw := network.BetweennessWeighted(g, path.DijkstraAllPaths(g))
			var nodes []struct {
				id    int64
				value float64
			}

			for k, v := range bw {
				nodes = append(nodes, struct {
					id    int64
					value float64
				}{k, v})
			}

			sort.SliceStable(nodes, func(i, j int) bool {
				return nodes[i].value < nodes[j].value
			})

			percentage := math.Round((10.0 / 100.0) * float64(len(nodes)))

			for i := 0; i < int(percentage); i++ {
				g.RemoveNode(nodes[i].id)

				new := g.NewNode()
				g.AddNode(new)

				edges := rand.Intn(d.MaxEdges)
				for i := 0; i < edges; i++ {
					edgeNode := rand.Intn(totalNodes)
					if new.ID() == int64(edgeNode) {
						continue
					}

					if g.Node(int64(edgeNode)) == nil {
						continue
					}

					edge := g.Edge(new.ID(), int64(edgeNode))
					if edge == nil {
						weight := rand.Intn(d.MaxWeight)
						g.SetWeightedEdge(&internal.AttrEdge{
							F:   new,
							T:   g.Node(int64(edgeNode)),
							W:   float64(weight),
							Dot: []encoding.Attribute{{Key: "label", Value: strconv.Itoa(weight)}},
						})
					}
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

			err = os.WriteFile(fmt.Sprintf("%v.dot", i), data, 0644)
			if err != nil {
				return nil, err
			}

			_, err = exec.Command("dot", "-Tsvg", fmt.Sprintf("%v.dot", i), "-o", fmt.Sprintf("%v.svg", i)).Output()
			// dot -Tsvg '/current/path/test.dot' -o output.svg
			if err != nil {
				return nil, err
			}
		}
	}

	return graphs, nil
}
