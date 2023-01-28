package calculate

import (
	"strconv"
	"testing"

	"github.com/bh90210/mlsic/internal"
	"github.com/stretchr/testify/assert"
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/encoding"
	"gonum.org/v1/gonum/graph/simple"
)

var graphNetworkTestCases = map[string]struct {
	numberOfNodes int
	opts          []NetworkOption
	nodeValue     map[NetworkOption]map[int64]float64
	mm            map[NetworkOption]MinMax
}{
	"empty graph": {
		nodeValue: map[NetworkOption]map[int64]float64{0: {}, 1: {}, 2: {}, 3: {}, 4: {}, 5: {}},
		mm:        map[NetworkOption]MinMax{0: {Min: 1.7976931348623157e+308, Max: 5e-324}, 1: {Min: 1.7976931348623157e+308, Max: 5e-324}, 2: {Min: 1.7976931348623157e+308, Max: 5e-324}, 3: {Min: 1.7976931348623157e+308, Max: 5e-324}, 4: {Min: 1.7976931348623157e+308, Max: 5e-324}, 5: {Min: 1.7976931348623157e+308, Max: 5e-324}},
	},
	"no options": {
		numberOfNodes: 10,
		nodeValue:     map[NetworkOption]map[int64]float64{0: {}, 1: {}, 2: {0: 0.1111111111111111, 1: 0.1111111111111111, 2: 0.1111111111111111, 3: 0.1111111111111111, 4: 0.1111111111111111, 5: 0.1111111111111111, 6: 0.1111111111111111, 7: 0.1111111111111111, 8: 0.1111111111111111, 9: 0.1111111111111111}, 3: {0: 9, 1: 9, 2: 9, 3: 9, 4: 9, 5: 9, 6: 9, 7: 9, 8: 9, 9: 9}, 4: {0: 9, 1: 9, 2: 9, 3: 9, 4: 9, 5: 9, 6: 9, 7: 9, 8: 9, 9: 9}, 5: {0: 4.5, 1: 4.5, 2: 4.5, 3: 4.5, 4: 4.5, 5: 4.5, 6: 4.5, 7: 4.5, 8: 4.5, 9: 4.5}},
		mm:            map[NetworkOption]MinMax{0: {Min: 1.7976931348623157e+308, Max: 5e-324}, 1: {Min: 1.7976931348623157e+308, Max: 5e-324}, 2: {Min: 0.1111111111111111, Max: 0.1111111111111111}, 3: {Min: 9, Max: 9}, 4: {Min: 9, Max: 9}, 5: {Min: 4.5, Max: 4.5}},
	},
	"with options": {
		numberOfNodes: 10,
		opts:          []NetworkOption{Farness},
		nodeValue:     map[NetworkOption]map[int64]float64{3: {0: 9, 1: 9, 2: 9, 3: 9, 4: 9, 5: 9, 6: 9, 7: 9, 8: 9, 9: 9}},
		mm:            map[NetworkOption]MinMax{3: {Min: 9, Max: 9}},
	},
}

func TestGraphNetwork(t *testing.T) {
	a := assert.New(t)

	for name, tc := range graphNetworkTestCases {
		t.Run(name, func(t *testing.T) {
			g := simple.NewWeightedUndirectedGraph(0, 0)
			for i := 0; i < tc.numberOfNodes; i++ {
				g.AddNode(g.NewNode())
			}

			for _, from := range graph.NodesOf(g.Nodes()) {
				for _, to := range graph.NodesOf(g.Nodes()) {
					if from.ID() == to.ID() {
						continue
					}

					edge := g.Edge(from.ID(), to.ID())
					if edge == nil {
						g.SetWeightedEdge(
							&internal.AttrEdge{
								F:   g.Node(from.ID()),
								T:   g.Node(to.ID()),
								W:   float64(1),
								Dot: []encoding.Attribute{{Key: "label", Value: strconv.Itoa(1)}},
							})
					}
				}
			}

			result := GraphNetwork(g, tc.opts...)
			a.Equal(tc.nodeValue, result.NodeValue)
			a.Equal(tc.mm, result.MM)
		})
	}
}
