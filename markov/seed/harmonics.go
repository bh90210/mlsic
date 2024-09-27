package seed

import (
	"time"

	"github.com/bh90210/mlsic"
	"github.com/bh90210/mlsic/markov"
)

// var _ mlsic.Harmonics = (*PrimeHarmonics)(nil)

// PrimeHarmonics .
type PrimeHarmonics struct {
	partials []mlsic.Partial
}

var prime = []mlsic.Partial{
	{
		Number:          2,
		AmplitudeFactor: .1,
		Start:           time.Duration(10 * time.Millisecond),
		Duration:        time.Duration(30 * time.Millisecond),
	},
	{
		Number:          3,
		AmplitudeFactor: .07,
		Start:           time.Duration(25 * time.Millisecond),
		Duration:        time.Duration(50 * time.Millisecond),
	},
	// {
	// 	Number:          4,
	// 	AmplitudeFactor: .6,
	// 	Start:           time.Duration(3 * time.Millisecond),
	// 	Duration:        time.Duration(1 * time.Millisecond),
	// },
	// {
	// 	Number:          5,
	// 	AmplitudeFactor: .5,
	// 	Start:           time.Duration(2 * time.Millisecond),
	// 	Duration:        time.Duration(2 * time.Millisecond),
	// },
	// {
	// 	Number:          6,
	// 	AmplitudeFactor: .5,
	// 	Start:           time.Duration(2 * time.Millisecond),
	// 	Duration:        time.Duration(2 * time.Millisecond),
	// },
	// {
	// 	Number:          7,
	// 	AmplitudeFactor: .2,
	// 	Start:           time.Duration(2 * time.Millisecond),
	// 	Duration:        time.Duration(2 * time.Millisecond),
	// },
	// {
	// 	Number:          8,
	// 	AmplitudeFactor: .15,
	// 	Start:           time.Duration(2 * time.Millisecond),
	// 	Duration:        time.Duration(2 * time.Millisecond),
	// },
	// {
	// 	Number:          9,
	// 	AmplitudeFactor: .1,
	// 	Start:           time.Duration(2 * time.Millisecond),
	// 	Duration:        time.Duration(2 * time.Millisecond),
	// },
	// {
	// 	Number:          10,
	// 	AmplitudeFactor: .05,
	// 	Start:           time.Duration(2 * time.Millisecond),
	// 	Duration:        time.Duration(2 * time.Millisecond),
	// },
}

// Partials .
func (p *PrimeHarmonics) Partials(poly []markov.Voice) []markov.Voice {
	// At this point partials only contain their number
	// and the amplitude factor.
	if p.partials == nil {
		// for i := 2; i < 230; i++ {
		// 	if big.NewInt(int64(i)).ProbablyPrime(0) {
		// 		p.partials = append(p.partials, mlsic.Partial{
		// 			Number:          i,
		// 			AmplitudeFactor: 1. / float64(i),
		// 		})
		// 	}
		// }
		for _, partial := range prime {
			p.partials = append(p.partials, partial)
		}
	}

	for i, voice := range poly {
		for toneIndex, tone := range voice {
			for _, partial := range p.partials {
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

			poly[i][toneIndex] = tone
		}
	}

	return poly
}

// Fundamental is always the first Wagon of a Train at index position zero.
const Fundamental = 0

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
