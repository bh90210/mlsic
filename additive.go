package mlsic

import (
	"errors"
	"math"
	"sort"
	"time"
)

// MaximumPartialStartingPoint .
const MaximumPartialStartingPoint = 1000

// MinimumPartialDuration .
const MinimumPartialDuration = 10

// ErrNotEnoughSpeakers .
var ErrNotEnoughSpeakers = errors.New("allowed number of speakers is 1+")

// Signal .
func Signal(poly []Voice, noOfSpeakers int) ([]Audio, error) {
	if noOfSpeakers < 1 {
		return nil, ErrNotEnoughSpeakers
	}
	signals := make([]Audio, noOfSpeakers)

	for _, voice := range poly {
		voiceIndex := voice.Ordered()
		voiceSignals := voice.Signals(noOfSpeakers)

		// We need to set the last phase of a sine
		// as the starting position of the next one.
		var previousFundamentalPhase float64
		// Range through the voice's tones.
		for _, i := range voiceIndex {
			// Set the tone to work this for this loop.
			tone := voice[i]
			// Generate fundamental's signal.
			phase, length, signal := tone.Signal(previousFundamentalPhase)

			// Set starting phase for next sine in voice.
			previousFundamentalPhase = phase

			// Create temporary slices for tone,
			toneSignal := make([][]float64, noOfSpeakers)
			// We the appropriate length for the duration of the fundamental.
			for o := range toneSignal {
				toneSignal[o] = make([]float64, length)
			}

			// Append fundamental's signal to the temporary buffer.
			for o, v := range signal {
				for speakerNumber := range noOfSpeakers {
					// Panning.
					panning := Panning(noOfSpeakers, speakerNumber, tone.Panning)

					toneSignal[speakerNumber][o] += v * tone.Fundamental.Amplitude * panning
				}
			}

			for _, partial := range tone.Partials {
				_, _, partialSignal := tone.PartialSignal(partial)
				for o, v := range partialSignal {
					for speakerNumber := range noOfSpeakers {
						// Panning.
						panning := Panning(noOfSpeakers, speakerNumber, tone.Panning)

						toneSignal[speakerNumber][o+partial.StartInSamples()] += v * (tone.Fundamental.Amplitude * partial.AmplitudeFactor) * panning
					}
				}
			}

			for speakerNo, signal := range toneSignal {
				for o, v := range signal {
					voiceSignals[speakerNo][i+o] = v
				}
			}
		}

		for speakerNo, signal := range voiceSignals {
			if len(signals[speakerNo]) < len(signal) {
				signals[speakerNo] = append(signals[speakerNo], make([]float64, len(signal)-len(signals[speakerNo]))...)
			}
		}

		for speakerNo, signal := range voiceSignals {
			for i, v := range signal {
				signals[speakerNo][i] += v
			}
		}
	}

	return signals, nil
}

// Voice is a single monophony from start to finish.
// It contains Tones, representing a fundamental
// (first index of the Tone) and it's partials.
type Voice map[int]Tone

// Ordered .
func (v Voice) Ordered() (voiceIndex []int) {
	voiceIndex = ordered(v)
	return
}

// Signals .
func (v Voice) Signals(noOfSpeakers int) (signals [][]float64) {
	// Determine the total trains length.
	var length int
	for k := range v {
		if length < k {
			length = k
		}
	}

	// Add the duration of voice's last tone.
	length += v[length].Fundamental.DurationInSamples()

	// Create signals slices of the appropriate length for each speaker.
	signals = make([][]float64, noOfSpeakers)
	for i := range signals {
		signals[i] = make([]float64, length+44100) // Add one extra second of silence at the end.
	}

	return
}

// LengthInSamples .
func (v Voice) LengthInSamples() (length int) {
	for k, v := range v {
		if k+v.Fundamental.DurationInSamples() > length {
			length = k + v.Fundamental.DurationInSamples()
		}
	}

	return
}

// Tone .
type Tone struct {
	Fundamental Sine
	Partials    []Partial
	// Panning information.
	Panning float64
}

// Signal creates a float64 audio signal out of the Sine and returns the length in samples.
// Note: signal always returns a signal to zero, or a full sine cycle.
// Inevitably it will return a sightly shorter signal than the original duration
// intended. This must be dealt with by the consumer.
func (t Tone) Signal(phase ...float64) (float64, int, Audio) {
	if phase != nil {
		t.Fundamental.phase = phase[0]
	}

	return signal(t.Fundamental.Frequency, t.Fundamental.phase, t.Fundamental.DurationInSamples())
}

// PartialSignal .
func (t Tone) PartialSignal(partial Partial) (float64, int, Audio) {
	frequency := t.Fundamental.Frequency * float64(partial.Number)
	if frequency > MaxFrequency {
		return 0, 0, nil
	}

	// TODO: should partials start at zero phase or follow fundamental's
	// at the particular point they start?
	return signal(frequency, .0, partial.DurationInSamples())
}

// Sine holds necessary data to construct a sine wave.
// It also has the method Signal() that creates the
// audio signal as Audio (float64 slice.)
type Sine struct {
	// Frequency of the sine wave.
	Frequency float64
	// Amplitude (velocity) of the sine wave.
	Amplitude float64
	// Duration of the sine wave in milliseconds.
	Duration time.Duration

	sampleFactor float64
	phase        float64
}

// DurationInSamples returns the assigned duration of Sine in samples.
func (s Sine) DurationInSamples() int {
	return SignalLengthMultiplier * int(s.Duration.Abs().Milliseconds())
}

func signal(frequency float64, phase float64, durationInSample int) (float64, int, Audio) {
	sampleFactor := frequency / SampleRate

	samples := make(Audio, durationInSample)
	for i := range samples {
		samples[i] = math.Sin(phase * 2.0 * math.Pi)
		_, phase = math.Modf(phase + sampleFactor)
	}

	return phase, len(samples), samples
}

func ordered[K int, V Tone](m map[K]V) (index []int) {
	for i := range m {
		index = append(index, int(i))
	}

	sort.Ints(index)

	return
}

// Partials .
func Partials(voice Voice, partials []Partial) Voice {
	for toneIndex, tone := range voice {
		for _, partial := range partials {
			if partial.Number*int(tone.Fundamental.Frequency) > MaxFrequency {
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

			tone.Partials = append(tone.Partials, Partial{
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
// 			if freq > MaxFrequency {
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

var PrimeMove1 = []Partial{
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
