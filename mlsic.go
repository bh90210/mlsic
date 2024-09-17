// Package mlsic .
package mlsic

import "time"

// SampleRate .
const SampleRate = 44100

// MaxFrequency .
const MaxFrequency = 18000

// SignalLengthMultiplier .
const SignalLengthMultiplier = 44

const (
	OneSpeaker         = 1
	TwoSpeakers        = 2
	SpeakerOne         = 0
	SpeakerTwo         = 1
	speakersMinPanning = 0
	speakersMaxPanning = 1
)

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
	Render(source []Audio, name string) error
}

// Sine .
type Sine struct {
	Frequency float64
	Amplitude float64
	Duration  time.Duration
}

// Scale .
func Scale(unscaledNum, minAllowed, maxAllowed, min, max float64) float64 {
	return (maxAllowed-minAllowed)*(unscaledNum-min)/(max-min) + minAllowed
}
