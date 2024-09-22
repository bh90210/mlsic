package seed

import (
	"math/rand/v2"
	"time"

	"github.com/bh90210/mlsic"
	"github.com/bh90210/mlsic/markov"
)

// PartialsGeneration .
func PartialsGeneration(h mlsic.Harmonics, voice mlsic.Voice) {
	partials := h.Partials()

	partialsTrains := make(mlsic.Voice)

	for trainIndex, train := range voice {
		wagon := train[0]
		partialsTrains[trainIndex] = make(mlsic.Train)

		for _, p := range partials {
			freq := wagon.Sine.Frequency * float64(p.Number)
			if freq > mlsic.MaxFrequency {
				continue
			}

			var partialIndex int
			for {
				// Double TODO: this needs to be less than the maximum duration of the fundamental.
				// TODO: this number determines the time in milliseconds the particular partial will start.
				// This is fundamental starting time + offset for the partial.
				// Number 500 meaning a partial can start up to half a second after the fundamental
				// is arbitrary. Fix it!
				partialIndex = rand.IntN(markov.MaximumPartialStartingPoint)
				if partialIndex != 0 {
					if _, ok := partialsTrains[trainIndex][partialIndex]; ok {
						continue
					}

					break
				}
			}

			l := mlsic.SignalLengthMultiplier * int(wagon.Sine.Duration.Abs().Milliseconds())
			l -= partialIndex
			l /= mlsic.SignalLengthMultiplier

			if l == 0 {
				l = markov.MinimumPartialDuration
			}

			partialsTrains[trainIndex][partialIndex] = mlsic.Wagon{
				Sine: mlsic.Sine{
					Frequency: freq,
					Amplitude: wagon.Sine.Amplitude * p.AmplitudeFactor,
					Duration:  time.Duration(l * int(time.Millisecond)),
				},
				// TODO: Panning of the partials is similar to fundamental. Make it dynamic.
				Panning: voice[trainIndex][0].Panning,
			}
		}
	}

	for trainIndex, train := range partialsTrains {
		for wagonIndex, wagon := range train {
			voice[trainIndex][wagonIndex] = wagon
		}
	}
}

func dynamicHarmonic(sine mlsic.Sine, partial int) float64 {
	f := mlsic.Scale((sine.Amplitude*sine.Frequency)*(float64(sine.Duration.Abs().Milliseconds())/float64(partial)), 0.0, 1.0, 0.0, 15000000.0)
	a := sine.Amplitude * f

	// return Scale(a, h.Partials1[partial], h.partials2[partial], 0.0, 1.0)
	return a
}
