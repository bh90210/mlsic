package seed

import (
	"math/big"
	"math/rand/v2"
	"time"

	"github.com/bh90210/mlsic"
	"github.com/bh90210/mlsic/markov"
)

// Harmonics .
type Harmonics struct {
	partials1 map[int]float64
	partials2 map[int]float64
}

func (h *Harmonics) gen() {
	h.partials1 = make(map[int]float64)

	for i := 2; i < 1000; i++ {
		h.partials1[i] = 1 / float64(i) / 10

		// If harmonic is odd mute.
		if i%2 == 1 {
			h.partials1[i] = 0.
		}

		if big.NewInt(int64(i)).ProbablyPrime(0) {
			h.partials1[i] = 0.
		}
	}

	h.partials2 = make(map[int]float64)

	for i := 2; i < 1000; i++ {
		h.partials2[i] = 0

		if big.NewInt(int64(i)).ProbablyPrime(0) {
			h.partials2[i] = 0.001 * float64(i)
		}
	}
}

// PartialsGeneration .
func (h *Harmonics) PartialsGeneration(trains markov.Trains) {
	if h.partials1 == nil || h.partials2 == nil {
		h.gen()
	}

	patrialsTrains := make(markov.Trains)

	for trainIndex, train := range trains {
		sine := train[0]
		patrialsTrains[trainIndex] = make(markov.Train)

		for partial, amplitude := range h.partials1 {
			freq := sine.Sine.Frequency * float64(partial)
			if freq > mlsic.MaxFrequency {
				continue
			}

			var partialIndex int
			for {
				// TODO: this number determines the time in milliseconds the particular partial will start.
				// This is fundamental starting time + offset for the partial.
				// Number 500 meaning a partial can start up to half a second after the fundamental
				// is arbitrary. Fix it!
				partialIndex = rand.IntN(markov.MaximumPartialStartingPoint)
				if partialIndex != 0 {
					if _, ok := patrialsTrains[trainIndex][partialIndex]; ok {
						continue
					}

					break
				}
			}

			l := mlsic.SignalLengthMultiplier * int(sine.Sine.Duration.Abs().Milliseconds())
			l -= partialIndex
			l /= mlsic.SignalLengthMultiplier

			if l == 0 {
				l = markov.MinimumPartialDuration
			}

			patrialsTrains[trainIndex][partialIndex] = markov.TrainContents{
				Sine: mlsic.Sine{
					Frequency: freq,
					Amplitude: sine.Sine.Amplitude * amplitude,
					Duration:  time.Duration(l * int(time.Millisecond)),
				},
				// TODO: Panning of the partials is similar to fundamental. Make it dynamic.
				Panning: trains[trainIndex][0].Panning,
			}
		}
	}

	for trainIndex, train := range patrialsTrains {
		for sineIndex, sine := range train {
			trains[trainIndex][sineIndex] = sine
		}
	}
}

func (h *Harmonics) dynamicHarmonic(sine mlsic.Sine, partial int) float64 {
	f := mlsic.Scale((sine.Amplitude*sine.Frequency)*(float64(sine.Duration.Abs().Milliseconds())/float64(partial)), 0.0, 1.0, 0.0, 15000000.0)
	a := sine.Amplitude * f

	// return mlsic.Scale(a, h.partials1[partial], h.partials2[partial], 0.0, 1.0)
	return a
}
