package seed

import (
	"math/big"
	"time"

	"github.com/bh90210/mlsic"
	"github.com/bh90210/mlsic/markov"
	"github.com/rs/zerolog/log"
)

// var _ mlsic.Harmonics = (*PrimeHarmonics)(nil)

// PrimeHarmonics .
type PrimeHarmonics struct {
	base []mlsic.Partial
}

func (p *PrimeHarmonics) init() {
	for i := 2; i < mlsic.MaxPartial; i++ {
		if big.NewInt(int64(i)).ProbablyPrime(0) {
			p.base = append(p.base, mlsic.Partial{
				Number:          i,
				AmplitudeFactor: 0.05,
				Start:           time.Duration(i) * time.Millisecond,
				Duration:        time.Duration(200) * time.Millisecond,
			})
		}
	}
}

// PartialsGen .
func (p *PrimeHarmonics) PartialsGen(voice markov.Voice) markov.Voice {
	if p.base == nil {
		p.init()
	}

	log.Debug().Any("p", p.base).Msg("yo")

	for toneIndex, tone := range voice {
		for _, partial := range p.base {
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
