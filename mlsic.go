// Package mlsic provides compositional algorithms
// and abstractions for producing audio.
package mlsic

import (
	"github.com/go-audio/audio"
	"gonum.org/v1/gonum/graph/simple"
)

// Graph is intended as decoupling the graph generation implementation from the compositional algorithms.
type Graph interface {
	// Dump must return a slice of simple.WeightedUndirectedGraph.
	// The order of the graphs are sorted in the slice is the order
	// they will be executed by the compositional algo.
	Dump() ([]*simple.WeightedUndirectedGraph, error)
}

// Renderer Pan declares a generic interface intended to be used for sound rendering operations.
type Renderer interface {
	// Render accepts a slice of audio.PCMBuffer that must always be of MONO format (ie. audio.FormatMono96000 or audio.FormatMono44100).
	// The number of buffers in the slice represent the numbers of audio channels the renderer must handle.
	Render(b []*audio.PCMBuffer) error
}

// Pan declares a generic interface intended to be used for panning operations.
// Pan is not designed for equal power multi-speaker panning. It will always
// pan the signal from a single source point.
type Pan interface {
	// Apply must accept the audio signal and the relevant position between total of channels Channels().
	// Buffer must always be of MONO format (ie. audio.FormatMono96000 or audio.FormatMono44100).
	// Position must be always between 0.0 and 1.0. In the case of two channels values from 0. to 0.5
	// will represent the left speaker and values from 0.5 to 1. the right one.
	// In the case of four speakers setup values from 0. to 0.25 will represent speakers no.1 etc..
	Apply(signal audio.PCMBuffer, position float32) []*audio.PCMBuffer
	// Channels must return the total number of channels the Apply() method should use.
	Channels() int
}
