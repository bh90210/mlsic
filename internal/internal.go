package internal

import (
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/encoding"
)

// DOTPrintOptions hold information relating to DOT files.
type DOTPrintOptions struct {
	DOTFile string
	SVGFile string
}

// AttrEdge is a custom edge with DOT attribute.
type AttrEdge struct {
	F, T graph.Node
	W    float64
	Dot  []encoding.Attribute
}

// From returns the from-node of the edge.
func (e *AttrEdge) From() graph.Node { return e.F }

// To returns the to-node of the edge.
func (e *AttrEdge) To() graph.Node { return e.T }

// ReversedEdge returns a new Edge with the F and T fields
// swapped.
func (e *AttrEdge) ReversedEdge() graph.Edge { e.F, e.T = e.T, e.F; return e }

// Weight returns the weight of the edge.
func (e *AttrEdge) Weight() float64 { return e.W }

// Attributes returns the DOT attributes of the edge.
func (e *AttrEdge) Attributes() []encoding.Attribute { return e.Dot }
