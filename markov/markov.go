package markov

import (
	"errors"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/bh90210/mlsic"
	"github.com/bh90210/mlsic/render"
	"github.com/go-audio/generator"
	"github.com/mb-14/gomarkov"
)

// Model .
type Model struct {
	Freq *gomarkov.Chain
	Amp  *gomarkov.Chain
	Dur  *gomarkov.Chain

	Poly *gomarkov.Chain
}

// Add .
func (m *Model) Add(train []Sine) {
	m.nilCheck()

	frequency := []string{}
	amplitude := []string{}
	duration := []string{}

	for _, v := range train {
		frequency = append(frequency, fmt.Sprintf("%f", v.Frequency))
		amplitude = append(amplitude, fmt.Sprintf("%f", v.Amplitude))
		duration = append(duration, fmt.Sprintf("%v", v.Duration.Milliseconds()))
	}

	m.Freq.Add(frequency)
	m.Amp.Add(amplitude)
	m.Dur.Add(duration)
}

type indexHelper struct {
	voice        int
	voiceIndex   int
	toneIndex    int
	partialIndex int
}

// AddPoly .
func (m *Model) AddPoly(poly []Voice) {
	m.nilCheck(true)

	// We will collect all indices there is a sine in a map[int].
	// the slice if ints []int are the voices this particular index
	// has a corresponding sine. This is because multiple voices may
	// share the same index.
	// For example: map[trainIndex+wagonIndex][voice][wagonIndex].
	indices := make(map[int][]indexHelper)
	for voiceNo, voice := range poly {
		for toneIndex, tone := range voice {
			indices[toneIndex] = append(indices[toneIndex], indexHelper{
				voice:        voiceNo,
				voiceIndex:   toneIndex,
				toneIndex:    toneIndex,
				partialIndex: -1,
			})

			for partialIndex, partial := range tone.Partials {
				indices[toneIndex+partial.StartInSamples()] = append(indices[toneIndex+partial.StartInSamples()], indexHelper{
					voice:        voiceNo,
					voiceIndex:   toneIndex,
					toneIndex:    toneIndex,
					partialIndex: partialIndex,
				})
			}
		}
	}

	// Order indices.
	OrderedIndices := make([]int, 0)
	for k := range indices {
		OrderedIndices = append(OrderedIndices, k)
	}

	sort.Ints(OrderedIndices)

	for _, i := range OrderedIndices {
		indexHelpers := indices[i]
		for _, h := range indexHelpers {
			tone := poly[h.voice][h.voiceIndex]
			t := []string{}

			// This means we are dealing with fundamental.
			switch h.partialIndex {
			case -1:
				if h.partialIndex == -1 {
					t = append(t, fmt.Sprintf("%f", tone.Fundamental.Frequency))
					t = append(t, fmt.Sprintf("%f", tone.Fundamental.Amplitude))
					t = append(t, fmt.Sprintf("%v", tone.Fundamental.DurationInSamples()))
					t = append(t, fmt.Sprintf("%f", tone.Panning))
					m.Poly.Add([]string{strings.Join(t, " ")})
				}

			default:
				partial := tone.Partials[h.partialIndex]

				t = append(t, fmt.Sprintf("%f", tone.Fundamental.Frequency*float64(partial.Number)))
				t = append(t, fmt.Sprintf("%f", tone.Fundamental.Amplitude*partial.AmplitudeFactor))
				t = append(t, fmt.Sprintf("%v", partial.DurationInSamples()))
				t = append(t, fmt.Sprintf("%f", tone.Panning))
			}

			m.Poly.Add([]string{strings.Join(t, " ")})
		}
	}
}

// Export .
func (m *Model) Export(path string) error {
	if m.Poly != nil {
		poly, err := m.Poly.MarshalJSON()
		if err != nil {
			return err
		}

		err = os.WriteFile(filepath.Join(path, "poly.json"), poly, 0644)
		if err != nil {
			return err
		}

		return nil
	}

	if m.Freq != nil {
		freq, err := m.Freq.MarshalJSON()
		if err != nil {
			return err
		}

		err = os.WriteFile(filepath.Join(path, "freq.json"), freq, 0644)
		if err != nil {
			return err
		}
	}

	if m.Amp != nil {
		amp, err := m.Amp.MarshalJSON()
		if err != nil {
			return err
		}

		err = os.WriteFile(filepath.Join(path, "amp.json"), amp, 0644)
		if err != nil {
			return err
		}
	}

	if m.Dur != nil {
		dur, err := m.Dur.MarshalJSON()
		if err != nil {
			return err
		}

		err = os.WriteFile(filepath.Join(path, "dur.json"), dur, 0644)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *Model) nilCheck(poly ...bool) {
	if poly != nil {
		if m.Poly == nil {
			m.Poly = gomarkov.NewChain(2)
		}

		return
	}

	if m.Freq == nil {
		m.Freq = gomarkov.NewChain(1)
	}

	if m.Amp == nil {
		m.Amp = gomarkov.NewChain(1)
	}

	if m.Dur == nil {
		m.Dur = gomarkov.NewChain(1)
	}
}

// Generate .
func Generate(filepath string, train []Sine, h mlsic.Harmonics, ngen int) {
	// Left channel.
	leftM := make(map[int][]float64, len(train))
	// Right channel.
	rightM := make(map[int][]float64, len(train))

	log.Info().Msg("generating train")

	var wg sync.WaitGroup
	var mu sync.Mutex

	partials := h.Partials()

	for i, v := range train {
		wg.Add(1)

		go func(i int, v Sine) {
			defer wg.Done()

			osc := generator.NewOsc(generator.WaveSine, v.Frequency, 44100)
			osc.Amplitude = v.Amplitude

			signal := osc.Signal(44 * int(v.Duration.Abs().Milliseconds()))

			for _, p := range partials {
				if v.Frequency*float64(p.Number) > 18000 {
					continue
				}

				osc = generator.NewOsc(generator.WaveSine, v.Frequency*float64(p.Number), 44100)

				if p.AmplitudeFactor < 0 {
					p.AmplitudeFactor *= -1
				}

				osc.Amplitude = v.Amplitude * p.AmplitudeFactor

				partialSignal := osc.Signal(44 * int(v.Duration.Abs().Milliseconds()))
				for i := range signal {
					signal[i] += partialSignal[i]
				}
			}

			for o, s := range signal {
				if s >= 1. {
					signal[o] = 0.
				} else if s <= -1. {
					signal[o] = 0.
				}
			}

			mu.Lock()
			leftM[i] = signal
			rightM[i] = signal
			mu.Unlock()
		}(i, v)
	}

	wg.Wait()

	// Left channel.
	var left []float64
	// Right channel.
	var right []float64

	for i := 0; i < len(leftM); i++ {
		left = append(left, leftM[i]...)
		right = append(right, rightM[i]...)
	}

	var music []mlsic.Audio
	music = append(music, mlsic.Audio(left))

	log.Info().Msg("rendering audio files")

	// Render.
	p := render.Wav{
		Filepath: filepath,
	}

	// p, err := render.NewPortAudio()
	// if err != nil {
	// 	log.Fatal().Err(err)
	// }

	if err := p.Render(music, fmt.Sprintf("ngen%v", ngen)); err != nil {
		log.Fatal().Err(err)
	}
}

// MaximumPartialStartingPoint .
const MaximumPartialStartingPoint = 1000

// MinimumPartialDuration .
const MinimumPartialDuration = 10

// ErrNotEnoughSpeakers .
var ErrNotEnoughSpeakers = errors.New("allowed number of speakers is 1+")

// Deconstruct .
func Deconstruct(poly []Voice, noOfSpeakers int) ([]mlsic.Audio, error) {
	if noOfSpeakers < 1 {
		return nil, ErrNotEnoughSpeakers
	}

	speakers := make([]mlsic.Audio, noOfSpeakers)

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
				for speakerNumber := 0; speakerNumber < noOfSpeakers; speakerNumber++ {
					// Panning.
					panning := mlsic.Panning(noOfSpeakers, speakerNumber, tone.Panning)

					toneSignal[speakerNumber][o] += v * tone.Fundamental.Amplitude * panning
				}
			}

			for _, partial := range tone.Partials {
				_, _, partialSignal := tone.PartialSignal(partial)
				for o, v := range partialSignal {
					for speakerNumber := 0; speakerNumber < noOfSpeakers; speakerNumber++ {
						// Panning.
						panning := mlsic.Panning(noOfSpeakers, speakerNumber, tone.Panning)

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
			if len(speakers[speakerNo]) < len(signal) {
				speakers[speakerNo] = append(speakers[speakerNo], make([]float64, len(signal)-len(speakers[speakerNo]))...)
			}
		}

		for speakerNo, signal := range voiceSignals {
			for i, v := range signal {
				speakers[speakerNo][i] += v
			}
		}
	}

	return speakers, nil
}

// Voice is a single monophony from start to finish.
// It contains Trains, representing a fundamental
// (first Wagon of the Train) and it's partials.
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

// Tone .
type Tone struct { // map[int]Partial
	Fundamental Sine
	Partials    []mlsic.Partial
	// Panning information.
	Panning float64
}

// Signal creates a float64 audio signal out of the Sine and returns the length in samples.
// Note: signal always returns a signal to zero, or a full sine cycle.
// Inevitably it will return a sightly shorter signal than the original duration
// intended. This must be dealt with by the consumer.
func (t Tone) Signal(phase ...float64) (float64, int, mlsic.Audio) {
	if phase != nil {
		t.Fundamental.phase = phase[0]
	}

	return signal(t.Fundamental.Frequency, t.Fundamental.phase, t.Fundamental.DurationInSamples())
}

// PartialSignal .
func (t Tone) PartialSignal(partial mlsic.Partial) (float64, int, mlsic.Audio) {
	frequency := t.Fundamental.Frequency * float64(partial.Number)
	if frequency > mlsic.MaxFrequency {
		return 0, 0, nil
	}

	// TODO: should partials start at zero phase or follow fundamental's
	// at the particular point they start?
	return signal(frequency, .0, partial.DurationInSamples())
}

// Sine holds necessary data to construct a sine wave.
// It also has the method Signal() that creates the
// audio signal as mlsic.Audio (float64 slice.)
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
	return mlsic.SignalLengthMultiplier * int(s.Duration.Abs().Milliseconds())
}

func signal(frequency float64, phase float64, durationInSample int) (float64, int, mlsic.Audio) {
	sampleFactor := frequency / mlsic.SampleRate

	samples := make(mlsic.Audio, durationInSample)
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
