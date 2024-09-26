package seed

import (
	"math/big"

	"github.com/bh90210/mlsic"
	"github.com/bh90210/mlsic/markov"
)

// var _ mlsic.Harmonics = (*PrimeHarmonics)(nil)

// PrimeHarmonics .
type PrimeHarmonics struct {
	partials []mlsic.Partial
}

// Partials .
func (p *PrimeHarmonics) Partials(poly []markov.Voice) {
	if len(p.partials) > 0 {
		return
	}

	for i := 2; i < 1000; i++ {
		v := 0.
		if big.NewInt(int64(i)).ProbablyPrime(0) {
			v = 0.0051 * float64(i)
		}

		p.partials = append(p.partials, mlsic.Partial{
			Number:          i,
			AmplitudeFactor: v,
		})
	}

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
