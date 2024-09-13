package seed

import (
	"math/big"
	"math/rand/v2"
	"sort"
	"time"

	"github.com/go-audio/generator"
	"github.com/rs/zerolog/log"

	"github.com/bh90210/mlsic"
	"github.com/bh90210/mlsic/markov"
)

// MaximumPartialStartingPoint .
const MaximumPartialStartingPoint = 1000

// MinimumPartialDuration .
const MinimumPartialDuration = 10

// MelodyTrain .
// func MelodyTrain() [][]mlsic.Sine {
func MelodyTrain() markov.Trains {
	log.Info().Msg("melody train")

	trains := make(markov.Trains, 0)

	// Manual .
	trains[0] = make(markov.Train)

	trains[0][0] = markov.TrainContents{
		Sine: mlsic.Sine{
			Frequency: 500,
			Amplitude: 0.49,
			Duration:  time.Duration(20 * time.Millisecond),
		},
		Panning: 0.5,
	}

	var prev int

	prev += int(trains[prev][0].Sine.Duration.Abs().Milliseconds()) *
		mlsic.SignalLengthMultiplier

	trains[prev] = make(markov.Train)

	trains[prev][0] = markov.TrainContents{
		Sine: mlsic.Sine{
			Frequency: 350,
			Amplitude: 0.4,
			Duration:  time.Duration(250 * time.Millisecond),
		},
		Panning: 0.5,
	}

	// prev += int(trains[prev][0].Duration.Abs().Milliseconds()) *
	// 	mlsic.SignalLengthMultiplier

	// trains[prev] = make(markov.Train)

	// trains[prev][0] = mlsic.Sine{
	// 	Frequency: 250,
	// 	Amplitude: 0.69,
	// 	Duration:  time.Duration(150 * time.Millisecond),
	// }

	// for i := 20; i > 0; i-- {
	// 	prev += int(trains[prev][0].Duration.Abs().Milliseconds()) *
	// 		mlsic.SignalLengthMultiplier

	// 	trains[prev] = make(markov.Train)

	// 	trains[prev][0] = mlsic.Sine{
	// 		Frequency: 2500 + float64(i),
	// 		Amplitude: 1. * float64(i) * 0.045,
	// 		Duration:  time.Duration(130. * time.Millisecond),
	// 	}
	// }

	// prev += int(trains[prev][0].Duration.Abs().Milliseconds()) *
	// 	mlsic.SignalLengthMultiplier

	// trains[prev] = make(markov.Train)

	// trains[prev][0] = mlsic.Sine{
	// 	Frequency: 500,
	// 	Amplitude: 0.5,
	// 	Duration:  time.Duration(180 * time.Millisecond),
	// }

	// for i := 99; i > 0; i-- {

	// 	prev += int(trains[prev][0].Duration.Abs().Milliseconds()) *
	// 		mlsic.SignalLengthMultiplier

	// 	trains[prev] = make(markov.Train)

	// 	f := 100 + float64(i)
	// 	a := float64(i) * 0.009

	// 	trains[prev][0] = mlsic.Sine{
	// 		Frequency: f,
	// 		Amplitude: a,
	// 		Duration:  time.Duration(24. * time.Millisecond),
	// 	}
	// }

	// // Random.
	// for i := 0; i < 5; i++ {
	// 	prev += int(trains[prev][0].Duration.Abs().Milliseconds()) *
	// 		mlsic.SignalLengthMultiplier

	// 	trains[prev] = make(markov.Train)

	// 	r := rand.Float64()
	// 	trains[prev][0] = mlsic.Sine{
	// 		Frequency: mlsic.Scale(r, 80, 500, 0., 1.),
	// 		Amplitude: mlsic.Scale(r, 0., 1., 0., 0.5) / 2,
	// 		Duration:  time.Duration(mlsic.Scale(r, 15, 150, 0., 1.)) * 2 * time.Millisecond,
	// 	}
	// }

	// Harmonics .
	var h harmonics
	h.partialsGeneration(trains)

	return trains
}

type harmonics struct {
	partials1 map[int]float64
	partials2 map[int]float64
}

func (h *harmonics) gen() {
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

func (h *harmonics) partialsGeneration(trains markov.Trains) {
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
				partialIndex = rand.IntN(MaximumPartialStartingPoint)
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
				l = MinimumPartialDuration
			}

			patrialsTrains[trainIndex][partialIndex] = markov.TrainContents{
				Sine: mlsic.Sine{
					Frequency: freq,
					Amplitude: sine.Sine.Amplitude * amplitude,
					Duration:  time.Duration(l * int(time.Millisecond)),
				},
			}
		}
	}

	for trainIndex, train := range patrialsTrains {
		for sineIndex, sine := range train {
			trains[trainIndex][sineIndex] = sine
		}
	}
}

func (h *harmonics) dynamicHarmonic(sine mlsic.Sine, partial int) float64 {
	f := mlsic.Scale((sine.Amplitude*sine.Frequency)*(float64(sine.Duration.Abs().Milliseconds())/float64(partial)), 0.0, 1.0, 0.0, 15000000.0)
	a := sine.Amplitude * f

	// return mlsic.Scale(a, h.partials1[partial], h.partials2[partial], 0.0, 1.0)
	return a
}

// DeconstructTrains .
func DeconstructTrains(trains markov.Trains) (mlsic.Audio, mlsic.Audio) {
	// Determine the total trains length.
	var length int
	for k := range trains {
		if length < k {
			length = k
		}
	}

	length += int(trains[length][0].Sine.Duration.Abs().Milliseconds() * mlsic.SignalLengthMultiplier)

	// Order trains map.
	trainKeys := make([]int, 0)
	for k := range trains {
		trainKeys = append(trainKeys, k)
	}

	sort.Ints(trainKeys)

	left := make([]float64, length+MaximumPartialStartingPoint+MinimumPartialDuration)
	right := make([]float64, length+MaximumPartialStartingPoint+MinimumPartialDuration)

	var previousSignalEnd int
	for _, i := range trainKeys {
		var signal []float64
		var fundamentalSignalEnd int

		for sineIndex, trainContent := range trains[i] {
			osc := generator.NewOsc(generator.WaveSine, trainContent.Sine.Frequency, mlsic.SampleRate)
			osc.Amplitude = trainContent.Sine.Amplitude

			sineSignal := osc.Signal(mlsic.SignalLengthMultiplier * int(trainContent.Sine.Duration.Abs().Milliseconds()))
			sineSignalEnd := trimToZero(sineSignal)

			// If we are dealing with the fundamental take note where it ends.
			if sineIndex == 0 {
				fundamentalSignalEnd = len(sineSignal) - sineSignalEnd
			}

			for p, s := range sineSignal {
				if s >= 1. {
					sineSignal[p] = 0.
				} else if s <= -1. {
					sineSignal[p] = 0.
				}
			}

			if len(signal) < len(sineSignal[:sineSignalEnd])+sineIndex {
				signal = append(signal, make([]float64, len(sineSignal[:sineSignalEnd])+sineIndex-len(signal))...)
			}

			for o, v := range sineSignal[:sineSignalEnd] {
				signal[sineIndex+o] += v
			}
		}

		for o, v := range signal {
			left[i+o-(i-previousSignalEnd)] += v
			right[i+o-(i-previousSignalEnd)] += v
		}

		previousSignalEnd += len(signal) - fundamentalSignalEnd
	}

	// Make sure the length of left/right channels is the same.
	if len(left) > len(right) {
		for range left[len(right):] {
			right = append(right, 0.)
		}
	} else if len(left) < len(right) {
		for range right[len(left):] {
			left = append(left, 0.)
		}
	}

	return mlsic.Audio(left), mlsic.Audio(right)
}

func trimToZero(s []float64) int {
	var lastIndex int

	switch {
	case s[len(s)-1] > 0:
		for i := len(s) - 1 - 1; i >= 0; i-- {
			if s[i] <= 0 {
				lastIndex = i
				break
			}
		}

	case s[len(s)-1] < 0:
		for i := len(s) - 1 - 1; i >= 0; i-- {
			if s[i] >= 0 {
				for o := i; o >= 0; o-- {
					if s[o] <= 0 {
						lastIndex = o
						break
					}
				}

				break
			}
		}

	case s[len(s)-1] == 0:
		lastIndex = len(s) - 1

	}

	return lastIndex
}
