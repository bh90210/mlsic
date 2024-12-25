package seed

import (
	"time"

	"github.com/bh90210/mlsic"
	"github.com/bh90210/mlsic/markov"
)

// Partials .
func Partials(voice markov.Voice, partials []mlsic.Partial) markov.Voice {
	for toneIndex, tone := range voice {
		for _, partial := range partials {
			if partial.Number*int(tone.Fundamental.Frequency) > mlsic.MaxFrequency {
				continue
			}

			var start time.Duration
			var duration time.Duration

			switch {
			// The duration of the tone is shorter than when partial begins.
			// Thus we skip this partial.
			case partial.Start > tone.Fundamental.Duration:
				continue

			// If the total duration of the partial is longer than the fundamental
			// then the partial needs to be shorter too.
			case partial.Start+partial.Duration > tone.Fundamental.Duration:
				start = partial.Start
				duration = tone.Fundamental.Duration - partial.Start

			default:
				start = partial.Start
				duration = partial.Duration
			}

			tone.Partials = append(tone.Partials, mlsic.Partial{
				Number:          partial.Number,
				AmplitudeFactor: partial.AmplitudeFactor,
				Start:           start,
				Duration:        duration,
			})
		}

		voice[toneIndex] = tone
	}

	return voice
}

// Fundamental is always the first Wagon of a Train at index position zero.
// const Fundamental = 0

// // PartialsGeneration .
// func PartialsGeneration(voice markov.Voice) {
// 	partialsTrains := make(markov.Voice)

// 	for toneIndex, tone := range voice {
// 		// At this point each tone has only one partial, the fundamental.
// 		fundamental := tone[Fundamental]
// 		partialsTrains[toneIndex] = make(markov.Tone)

// 		// Init prime harmonics with the fundamental.
// 		fundamentalHarmonics := PrimeHarmonics{
// 			Fundamental: &fundamental,
// 		}

// 		// Generate the partials.
// 		partials := fundamentalHarmonics.Partials()

// 		// Range through them and append them to the tone.
// 		for _, partial := range partials {
// 			freq := fundamental.Sine.Frequency * float64(partial.Number)
// 			if freq > mlsic.MaxFrequency {
// 				continue
// 			}

// 			partialsTrains[toneIndex][partial.StartInSamples()] = markov.Partial{
// 				Sine: markov.Sine{
// 					Frequency: freq,
// 					Amplitude: fundamental.Sine.Amplitude * partial.AmplitudeFactor,
// 					Duration:  partial.Duration,
// 				},
// 				// TODO: Panning of the partials is similar to fundamental. Make it dynamic.
// 				Panning:        voice[toneIndex][Fundamental].Panning,
// 				NotFundamental: &partial,
// 			}
// 		}
// 	}

// 	for toneIndex, tone := range partialsTrains {
// 		for partialIndex, partial := range tone {
// 			voice[toneIndex][partialIndex] = partial
// 		}
// 	}
// }

var PrimeMove1 = []mlsic.Partial{
	{
		Number:          2,
		AmplitudeFactor: .01,
		Start:           time.Duration(10 * time.Millisecond),
		Duration:        time.Duration(300 * time.Millisecond),
	},
	{
		Number:          3,
		AmplitudeFactor: .007,
		Start:           time.Duration(25 * time.Millisecond),
		Duration:        time.Duration(500 * time.Millisecond),
	},
	{
		// Number:          5,
		Number:          4,
		AmplitudeFactor: .06,
		Start:           time.Duration(3 * time.Millisecond),
		Duration:        time.Duration(1000 * time.Millisecond),
	},
	{
		// Number:          7,
		Number:          5,
		AmplitudeFactor: .05,
		Start:           time.Duration(50 * time.Millisecond),
		Duration:        time.Duration(200 * time.Millisecond),
	},
	{
		// Number:          11,
		Number:          6,
		AmplitudeFactor: .05,
		Start:           time.Duration(30 * time.Millisecond),
		Duration:        time.Duration(300 * time.Millisecond),
	},
	{
		// Number:          13,
		Number:          7,
		AmplitudeFactor: .02,
		Start:           time.Duration(35 * time.Millisecond),
		Duration:        time.Duration(2000 * time.Millisecond),
	},
	{
		// Number:          17,
		Number:          8,
		AmplitudeFactor: .015,
		Start:           time.Duration(100 * time.Millisecond),
		Duration:        time.Duration(1000 * time.Millisecond),
	},
	{
		// Number:          19,
		Number:          9,
		AmplitudeFactor: .001,
		Start:           time.Duration(200 * time.Millisecond),
		Duration:        time.Duration(2000 * time.Millisecond),
	},
	{
		// Number:          23,
		Number:          10,
		AmplitudeFactor: .005,
		Start:           time.Duration(2 * time.Millisecond),
		Duration:        time.Duration(200 * time.Millisecond),
	},
}
