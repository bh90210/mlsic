package markov

import (
	"fmt"
	"os"
	"path/filepath"
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
}

// Add .
func (t *Models) Add(train []mlsic.Sine) {
	t.nilCheck()

	frequency := []string{}
	amplitude := []string{}
	duration := []string{}

	for _, v := range train {
		frequency = append(frequency, fmt.Sprintf("%f", v.Frequency))
		amplitude = append(amplitude, fmt.Sprintf("%f", v.Amplitude))
		duration = append(duration, fmt.Sprintf("%v", v.Duration.Milliseconds()))
	}

	t.Freq.Add(frequency)
	t.Amp.Add(amplitude)
	t.Dur.Add(duration)
}

// Export .
func (t *Models) Export(path string) error {
	t.nilCheck()

	freq, err := t.Freq.MarshalJSON()
	if err != nil {
		return err
	}

	err = os.WriteFile(filepath.Join(path, "freq.json"), freq, 0644)
	if err != nil {
		return err
	}

	amp, err := t.Amp.MarshalJSON()
	if err != nil {
		return err
	}

	err = os.WriteFile(filepath.Join(path, "amp.json"), amp, 0644)
	if err != nil {
		return err
	}

	dur, err := t.Dur.MarshalJSON()
	if err != nil {
		return err
	}

	err = os.WriteFile(filepath.Join(path, "dur.json"), dur, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (t *Models) nilCheck() {
	if t.Freq == nil {
		t.Freq = &gomarkov.Chain{}
	}

	if t.Amp == nil {
		t.Amp = &gomarkov.Chain{}
	}

	if t.Dur == nil {
		t.Dur = &gomarkov.Chain{}
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

			for partial, amplitude := range h.Partials {
				if v.Frequency*float64(partial) > 18000 {
					continue
				}

				osc = generator.NewOsc(generator.WaveSine, v.Frequency*float64(partial), 44100)

				amplitude = amplitude + ((amplitude - h.Partials[partial]) / 2)
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

// Train fundamental and partials.
type Train map[int]TrainContents

// Trains is a single voice (monophony) from start to finish.
type Trains map[int]Train

// Poly each slice is a different poly voice.
type Poly []Trains

// TrainContents .
type TrainContents struct {
	Sine    mlsic.Sine
	Panning float64
}
