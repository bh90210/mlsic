package markov

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"github.com/rs/zerolog/log"

	"github.com/bh90210/mlsic"
	"github.com/bh90210/mlsic/render"
	"github.com/go-audio/generator"
	"github.com/mb-14/gomarkov"
)

// Models .
type Models struct {
	Freq *gomarkov.Chain
	Amp  *gomarkov.Chain
	Dur  *gomarkov.Chain

	Poly *gomarkov.Chain
}

// Add .
func (m *Models) Add(train []mlsic.Sine) {
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
	voice      int
	trainIndex int
	wagonIndex int
}

// AddPoly .
func (m *Models) AddPoly(poly Poly) {
	m.nilCheck(true)

	// We will collect all indices there is a sine in a map[int].
	// the slice if ints []int are the voices this particular index
	// has a corresponding sine. This is because multiple voices may
	// share the same index.
	// For example: map[trainIndex+wagonIndex][voice][wagonIndex].
	indices := make(map[int][]indexHelper)
	for voiceNo, voice := range poly {
		for trainIndex, train := range voice {
			for wagonIndex := range train {
				indices[trainIndex+wagonIndex] = append(indices[trainIndex+wagonIndex], indexHelper{
					voice:      voiceNo,
					trainIndex: trainIndex,
					wagonIndex: wagonIndex,
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
			wagon := poly[h.voice][h.trainIndex][h.wagonIndex]

			w := []string{}
			w = append(w, fmt.Sprintf("%f", wagon.Sine.Frequency))
			w = append(w, fmt.Sprintf("%f", wagon.Sine.Amplitude))
			w = append(w, fmt.Sprintf("%v", wagon.Sine.Duration.Milliseconds()))
			w = append(w, fmt.Sprintf("%f", wagon.Panning))

			m.Poly.Add([]string{strings.Join(w, " ")})
		}
	}
}

// Export .
func (m *Models) Export(path string) error {
	if m.Poly != nil {
		poly, err := m.Poly.MarshalJSON()
		if err != nil {
			return err
		}

		err = os.WriteFile(filepath.Join(path, "poly.json"), poly, 0644)
		if err != nil {
			return err
		}
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

func (m *Models) nilCheck(poly ...bool) {
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
func Generate(filepath string, train []mlsic.Sine, h *Harmonics, ngen int) {
	// Left channel.
	leftM := make(map[int][]float64, len(train))
	// Right channel.
	rightM := make(map[int][]float64, len(train))

	log.Info().Msg("generating train")

	var wg sync.WaitGroup
	var mu sync.Mutex

	for i, v := range train {
		wg.Add(1)

		go func(i int, v mlsic.Sine) {
			defer wg.Done()

			osc := generator.NewOsc(generator.WaveSine, v.Frequency, 44100)
			osc.Amplitude = v.Amplitude

			signal := osc.Signal(44 * int(v.Duration.Abs().Milliseconds()))

			for partial, amplitude := range h.Partials1 {
				if v.Frequency*float64(partial) > 18000 {
					continue
				}

				osc = generator.NewOsc(generator.WaveSine, v.Frequency*float64(partial), 44100)

				amplitude = amplitude + ((amplitude - h.Partials1[partial]) / 2)
				if amplitude < 0 {
					amplitude *= -1
				}

				osc.Amplitude = v.Amplitude * amplitude

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

// Poly each slice is a different poly voice.
type Poly []Voice

// Voice is a single voice (monophony) from start to finish.
type Voice map[int]Train

// Train fundamental and partials.
type Train map[int]Wagon

// Wagon .
type Wagon struct {
	Sine    mlsic.Sine
	Panning float64
}

// MaximumPartialStartingPoint .
const MaximumPartialStartingPoint = 1000

// MinimumPartialDuration .
const MinimumPartialDuration = 10

// ErrNotEnoughSpeakers .
var ErrNotEnoughSpeakers = errors.New("allowed number of speakers is 1+")

// DeconstructTrains .
func DeconstructTrains(poly Poly, noOfSpeakers int) ([]mlsic.Audio, error) {
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

				// Append empty values to signal if signal is shorter than needed.
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

		for o, s := range speakersSignals {
			if len(speakers[o]) < len(s) {
				speakers[o] = append(speakers[o], make([]float64, len(s)-len(speakers[o]))...)
			}
		}

		for o, s := range speakersSignals {
			for p, v := range s {
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
