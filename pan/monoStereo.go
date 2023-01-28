// Package pan holds various mlsic.Pan implementations.
// It should be used along top package compositional algos.
package pan

import (
	"github.com/bh90210/mlsic"
	"github.com/go-audio/audio"
)

var _ mlsic.Pan = (*Mono)(nil)

// Mono implements mlsic.Pan. It accepts no options.
type Mono struct{}

// Apply will append the signal to an audio.PCMBuffer slice and ignore the position.
func (m *Mono) Apply(signal audio.PCMBuffer, position float32) []*audio.PCMBuffer {
	var monoBuffer []*audio.PCMBuffer
	monoBuffer = append(monoBuffer, &signal)

	return monoBuffer
}

// Channels will return 1 (mono).
func (m *Mono) Channels() int {
	return 1
}

var _ mlsic.Pan = (*LinearStereo)(nil)

// LinearStereo implements mlsic.Pan. It accepts no options.
type LinearStereo struct{}

// Apply Apply will append the signal to an audio.PCMBuffer slice
// and apply linear power panning to it determined by position.
func (l *LinearStereo) Apply(signal audio.PCMBuffer, position float32) []*audio.PCMBuffer {
	var pannedBuffers []*audio.PCMBuffer
	for i := 0; i < 2; i++ {
		pannedBuffers = append(pannedBuffers,
			&audio.PCMBuffer{
				Format: signal.Format,
			})
	}

	for i, b := range pannedBuffers {
		switch i {
		case 0:
			var leftSignal []float32
			for _, v := range signal.F32 {
				leftSignal = append(leftSignal, v*(1-position))
			}

			b.F32 = append(b.F32, leftSignal...)

		case 1:
			var rightSignal []float32
			for _, v := range signal.F32 {
				rightSignal = append(rightSignal, v*position)
			}

			b.F32 = append(b.F32, rightSignal...)
		}
	}

	return pannedBuffers
}

// Channels will return 2 (stereo).
func (l *LinearStereo) Channels() int {
	return 2
}
