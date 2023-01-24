// Package calculate provides convenient functions to help perform various
// repetitive calculations. It is meant to be used by the top level package
// inside compositional algos (ie. Algo1.)
package calculate

import (
	"math"
	"sync"

	"gonum.org/v1/gonum/graph/network"
	"gonum.org/v1/gonum/graph/path"
	"gonum.org/v1/gonum/graph/simple"
)

const (
	// Betweenness centrality is a measure of centrality in a graph based on shortest paths.
	Betweenness NetworkOption = iota
	// BetweennessWeighted centrality is a measure of centrality in a graph based on shortest paths.
	BetweennessWeighted
	// Closeness is a measure of centrality in a network, calculated as the reciprocal of the sum of the length of the shortest paths between the node and all other nodes in the graph.
	Closeness
	// Farness is the opposite of Closeness.
	Farness
	// Harmonic Closeness.
	Harmonic
	// Residual Closeness.
	Residual
)

// NetworkOption is a custom type around the various network calculations we need to perform.
type NetworkOption int

// MinMax is a structure that holds the minimum and maximum of something.
// It is used to hold the min/max range of the various calculations done to graphs
// as it is needed when we need to scale things (see calculate.LinearScale().)
type MinMax struct {
	Min float64
	Max float64
}

// NetworkAnalysis
type NetworkAnalysis struct {
	NodeValue map[NetworkOption]map[int64]float64
	MM        map[NetworkOption]MinMax

	graph       *simple.WeightedUndirectedGraph
	allShortest path.AllShortest

	wg sync.WaitGroup
	mu sync.Mutex
}

// GraphNetwork .
func GraphNetwork(g *simple.WeightedUndirectedGraph, opts ...NetworkOption) *NetworkAnalysis {
	ga := &NetworkAnalysis{
		NodeValue: make(map[NetworkOption]map[int64]float64),
		MM:        make(map[NetworkOption]MinMax),

		graph:       g,
		allShortest: path.DijkstraAllPaths(g),
	}

	if opts == nil {
		opts = append(opts, []NetworkOption{Betweenness, BetweennessWeighted, Closeness, Farness, Harmonic, Residual}...)
	}

	ga.wg.Add(len(opts))

	for _, opt := range opts {
		go ga.calc(opt)
	}

	ga.wg.Wait()

	return ga
}

func (ga *NetworkAnalysis) calc(opt NetworkOption) {
	defer ga.wg.Done()

	var r map[int64]float64

	switch opt {
	case Betweenness:
		r = network.Betweenness(ga.graph)
	case BetweennessWeighted:
		r = network.BetweennessWeighted(ga.graph, ga.allShortest)
	case Closeness:
		r = network.Closeness(ga.graph, ga.allShortest)
	case Farness:
		r = network.Farness(ga.graph, ga.allShortest)
	case Harmonic:
		r = network.Harmonic(ga.graph, ga.allShortest)
	case Residual:
		r = network.Residual(ga.graph, ga.allShortest)
	}

	min := math.MaxFloat64
	max := math.SmallestNonzeroFloat64

	for _, v := range r {
		if v < min {
			min = v
		}

		if v > max {
			max = v
		}
	}

	ga.mu.Lock()
	ga.NodeValue[opt] = r
	ga.MM[opt] = MinMax{Min: min, Max: max}
	ga.mu.Unlock()
}
