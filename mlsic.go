// Package mlsic .
package mlsic

// Audio is a 64 bit float slice with PCM signal values from -1.0 to 1.0.
type Audio []float64

// Reader .
type Reader interface {
	Read(Audio) (int, error)
}

// Writer .
type Writer interface {
	Write(Audio) (int, error)
}

// Renderer .
type Renderer interface {
	Render([]Audio) error
}
