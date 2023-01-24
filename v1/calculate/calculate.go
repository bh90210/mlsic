// Package calculate provides convenient functions to help perform various
// repetitive calculations. It is meant to be used by the top level package
// inside compositional algos (ie. Algo1.)
package calculate

import (
	"github.com/chewxy/math32"
	"github.com/go-audio/audio"
)

// LinearScale (value - r.Min) / (r.Max - r.Min) * (t.Max - t.Min) + t.Min.
func LinearScale(value float64, r MinMax, t MinMax) float64 {
	return ((value - r.Min) / (r.Max - r.Min)) * ((t.Max - t.Min) + t.Min)
}

// SineOptions hold the frequency, amplitude, and attack, decay, release information.
type SineOptions struct {
	// Freq is the frequency of the sine wave.
	Freq float32
	// Amp is the amplitude of the sine wave.
	Amp float32
	// A, D, R represents the attack, decay, release times in milliseconds.
	A, D, R float64
}

// SineGeneration generates a sine wave according to the parameters provided.
// It returns an audio.PCMBuffer loaded with the generated sine wave as 32bit floating point slice []F32.
// The math package used to calculate the sine wave is github.com/chewxy/math32.
// If f is nil it returns an empty audio.PCMbuffer{}.
func SineGeneration(f *audio.Format, opts SineOptions) audio.PCMBuffer {
	if f == nil {
		return audio.PCMBuffer{}
	}

	var phase float32
	var step = opts.Freq / float32(f.SampleRate)
	var signal = audio.PCMBuffer{
		Format: f,
	}

	attackInSamples := int(opts.A / (1000 / float64(f.SampleRate)))
	for i := 0; i < attackInSamples; i++ {
		sample := math32.Sin(2 * math32.Pi * phase)

		s := sample * (float32(i) * (opts.Amp / float32(attackInSamples)))

		signal.F32 = append(signal.F32, s)

		_, phase = math32.Modf(phase + step)
	}

	decayInSamples := int(opts.D / (1000 / float64(f.SampleRate)))
	for i := 0; i < decayInSamples; i++ {
		sample := math32.Sin(2 * math32.Pi * phase)

		s := sample * opts.Amp

		signal.F32 = append(signal.F32, s)

		_, phase = math32.Modf(phase + step)
	}

	releaseInSamples := int(opts.R / (1000 / float64(f.SampleRate)))
	for i := releaseInSamples; i > 0; i-- {
		sample := math32.Sin(2 * math32.Pi * phase)

		s := sample * (float32(i) * (opts.Amp / float32(releaseInSamples)))

		signal.F32 = append(signal.F32, s)

		_, phase = math32.Modf(phase + step)
	}

	return signal
}
