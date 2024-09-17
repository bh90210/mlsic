package seed

import (
	"errors"
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
func MelodyTrain() markov.Poly {
	log.Info().Msg("melody train")

	voice1Trains := make(markov.Trains, 0)

	// Manual .
	voice1Trains[0] = make(markov.Train)

	voice1Trains[0][0] = markov.TrainContents{
		Sine: mlsic.Sine{
			Frequency: 555,
			Amplitude: 0.49,
			Duration:  time.Duration(400 * time.Millisecond),
		},
		Panning: 0.,
	}

	var prev int

	prev += int(voice1Trains[prev][0].Sine.Duration.Abs().Milliseconds()) *
		mlsic.SignalLengthMultiplier

	voice1Trains[prev] = make(markov.Train)

	voice1Trains[prev][0] = markov.TrainContents{
		Sine: mlsic.Sine{
			Frequency: 350,
			Amplitude: 0.4,
			Duration:  time.Duration(400 * time.Millisecond),
		},
		Panning: 0.25,
	}

	// Harmonics .
	var h harmonics
	h.partialsGeneration(voice1Trains)

	var poly markov.Poly
	poly = append(poly, voice1Trains)

	//
	// Second voice.
	//

	voice2Trains := make(markov.Trains, 0)

	// Manual .
	voice2Trains[0] = make(markov.Train)

	voice2Trains[0][0] = markov.TrainContents{
		Sine: mlsic.Sine{
			Frequency: 1000,
			Amplitude: 0.29,
			Duration:  time.Duration(180 * time.Millisecond),
		},
		Panning: 0.50,
	}

	prev = 0

	prev += int(voice2Trains[prev][0].Sine.Duration.Abs().Milliseconds()) *
		mlsic.SignalLengthMultiplier

	voice2Trains[prev] = make(markov.Train)

	voice2Trains[prev][0] = markov.TrainContents{
		Sine: mlsic.Sine{
			Frequency: 750,
			Amplitude: 0.2,
			Duration:  time.Duration(250 * time.Millisecond),
		},
		Panning: 0.75,
	}

	// Harmonics .
	h.partialsGeneration(voice2Trains)

	poly = append(poly, voice2Trains)

	return poly
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

func (h *harmonics) dynamicHarmonic(sine mlsic.Sine, partial int) float64 {
	f := mlsic.Scale((sine.Amplitude*sine.Frequency)*(float64(sine.Duration.Abs().Milliseconds())/float64(partial)), 0.0, 1.0, 0.0, 15000000.0)
	a := sine.Amplitude * f

	// return mlsic.Scale(a, h.partials1[partial], h.partials2[partial], 0.0, 1.0)
	return a
}

// ErrNotEnoughSpeakers .
var ErrNotEnoughSpeakers = errors.New("allowed number of speakers is 1+")

// DeconstructTrains .
func DeconstructTrains(poly markov.Poly, noOfSpeakers int) ([]mlsic.Audio, error) {
	if noOfSpeakers < 1 {
		return nil, ErrNotEnoughSpeakers
	}

	speakers := make([]mlsic.Audio, noOfSpeakers)

	for _, voiceTrains := range poly {
		// Determine the total trains length.
		var length int
		for k := range voiceTrains {
			if length < k {
				length = k
			}
		}

		length += int(voiceTrains[length][0].Sine.Duration.Abs().Milliseconds() * mlsic.SignalLengthMultiplier)

		// Order trains map.
		trainKeys := make([]int, 0)
		for k := range voiceTrains {
			trainKeys = append(trainKeys, k)
		}

		sort.Ints(trainKeys)

		speakersSignals := make([][]float64, noOfSpeakers)
		for i := range speakersSignals {
			speakersSignals[i] = make([]float64, length+MaximumPartialStartingPoint+MinimumPartialDuration)
		}

		var previousSignalEnd int
		for _, i := range trainKeys {
			signal := make([][]float64, noOfSpeakers)

			var fundamentalSignalEnd int
			for sineIndex, trainContent := range voiceTrains[i] {
				osc := generator.NewOsc(generator.WaveSine, trainContent.Sine.Frequency, mlsic.SampleRate)
				osc.Amplitude = trainContent.Sine.Amplitude

				sineSignal := osc.Signal(mlsic.SignalLengthMultiplier * int(trainContent.Sine.Duration.Abs().Milliseconds()))
				sineSignalEnd := trimToZero(sineSignal)

				// If we are dealing with the fundamental take note where it ends.
				if sineIndex == 0 {
					fundamentalSignalEnd = len(sineSignal) - sineSignalEnd
				}

				// TODO: This is crude, fix it.
				// Prevent audio clipping.
				for p, s := range sineSignal {
					if s >= 1. {
						sineSignal[p] = 0.
					} else if s <= -1. {
						sineSignal[p] = 0.
					}
				}

				// Append empty values to signal if signal is shorted than needed.
				for o := 0; o < noOfSpeakers; o++ {
					if len(signal[o]) < len(sineSignal[:sineSignalEnd])+sineIndex {
						signal[o] = append(signal[o], make([]float64, len(sineSignal[:sineSignalEnd])+sineIndex-len(signal[o]))...)
					}
				}

				for o, v := range sineSignal[:sineSignalEnd] {
					// Panning.
					for speakerNumber := 0; speakerNumber < noOfSpeakers; speakerNumber++ {
						var panning float64

						switch noOfSpeakers {
						// Mono.
						case mlsic.OneSpeaker:
							panning = 1

						// Stereo.
						case mlsic.TwoSpeakers:
							switch speakerNumber {
							// Left.
							case mlsic.SpeakerOne:
								panning = 1 - trainContent.Panning

							// Right.
							case mlsic.SpeakerTwo:
								panning = trainContent.Panning
							}

						// Three and more speakers.
						default:
							// Find the width of individual speaker.
							speakerWidth := 1. / float64(noOfSpeakers)
							// Find the max width value of current speaker.
							speakerMax := speakerWidth * float64((speakerNumber))
							// Find the min width value of current speaker.
							speakerMin := speakerMax - speakerWidth
							// Find current speaker's mid point.
							speakerMid := speakerMin + (speakerWidth / 2)

							switch {
							// If the panning value is within the width of this speaker and
							// above or below speaker's mid.
							case trainContent.Panning >= speakerMin &&
								trainContent.Panning <= speakerMax:

								if trainContent.Panning == speakerMid {
									panning = 1
								}

								if trainContent.Panning < speakerMid {
									panning = 1 - mlsic.Scale(speakerMid-trainContent.Panning, 0., 1., 0., speakerWidth)
								}

								if trainContent.Panning > speakerMid {
									panning = mlsic.Scale(trainContent.Panning-speakerMid, 0., 1., 0., speakerWidth)
								}

							// If panning value is above this speaker's range.
							// This implies that there is a speaker on the right.
							case trainContent.Panning > speakerMax &&
								trainContent.Panning < (speakerMid+speakerWidth):
								panning = mlsic.Scale(speakerMid+speakerWidth-trainContent.Panning, 0., 1., 0., speakerWidth)

							// If panning value is bellow this speaker's range.
							// This implies that there is a speaker on the left.
							case trainContent.Panning < speakerMin &&
								trainContent.Panning > (speakerMid-speakerWidth):
								panning = 1 - mlsic.Scale(speakerMid-trainContent.Panning, 0., 1., 0., speakerWidth)

							case speakerNumber == 0 &&
								trainContent.Panning > speakerMid+(speakerWidth*float64(noOfSpeakers-1)):
								panning = mlsic.Scale(trainContent.Panning-(speakerMid+(speakerWidth*float64(noOfSpeakers-1))), 0., 1., 0., speakerWidth)

							case speakerNumber == noOfSpeakers &&
								trainContent.Panning < speakerWidth/2:
								panning = mlsic.Scale((speakerWidth/2)-trainContent.Panning, 0., 1., 0., speakerWidth)
							}
						}

						signal[speakerNumber][sineIndex+o] += v * panning
					}
				}
			}

			for p, s := range signal {
				for o, v := range s {
					speakersSignals[p][i+o-(i-previousSignalEnd)] += v
				}
			}

			previousSignalEnd += len(signal[0]) - fundamentalSignalEnd
		}

		for o, cs := range speakersSignals {
			if len(speakers[o]) < len(cs) {
				speakers[o] = append(speakers[o], make([]float64, len(cs)-len(speakers[o]))...)
			}
		}

		for o, cs := range speakersSignals {
			for p, v := range cs {
				speakers[o][p] += v
			}
		}
	}

	return speakers, nil
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
