package markov

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"

	"github.com/bh90210/mlsic"
	"github.com/bh90210/mlsic/render"
	"github.com/go-audio/generator"
	"github.com/mb-14/gomarkov"
)

// Train .
type Train struct {
	Freq *gomarkov.Chain
	Amp  *gomarkov.Chain
	Dur  *gomarkov.Chain
}

// Add .
func (m *Train) Add(train []mlsic.Sine) {
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

// Export .
func (m *Train) Export(path string) error {
	freq, err := m.Freq.MarshalJSON()
	if err != nil {
		return err
	}

	err = os.WriteFile(filepath.Join(path, "freq.json"), freq, 0644)
	if err != nil {
		return err
	}

	amp, err := m.Amp.MarshalJSON()
	if err != nil {
		return err
	}

	err = os.WriteFile(filepath.Join(path, "amp.json"), amp, 0644)
	if err != nil {
		return err
	}

	dur, err := m.Dur.MarshalJSON()
	if err != nil {
		return err
	}

	err = os.WriteFile(filepath.Join(path, "dur.json"), dur, 0644)
	if err != nil {
		return err
	}

	return nil
}

// Generate .
func Generate(filepath string, train []mlsic.Sine, h *Harmonics) {
	// Left channel.
	var left []float64
	// Right channel.
	var right []float64
	for _, v := range train {
		osc := generator.NewOsc(generator.WaveSine, v.Frequency, 44100)
		osc.Amplitude = v.Amplitude
		// osc.SetAttackInMs(10)

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
			// osc.SetAttackInMs(10)

			partialSignal := osc.Signal(44 * int(v.Duration.Abs().Milliseconds()))
			for i := range signal {
				signal[i] += partialSignal[i]
			}
		}

		for i, s := range signal {
			if s >= 1. {
				signal[i] = 0.
			} else if s <= -1. {
				signal[i] = 0.
			}
		}

		left = append(left, signal...)
		right = append(right, signal...)
	}

	var music []mlsic.Audio
	music = append(music, mlsic.Audio(left), mlsic.Audio(right))

	// Render.
	p := render.Wav{
		Filepath: filepath,
	}

	// p, err := render.NewPortAudio()
	// if err != nil {
	// 	log.Fatal().Err(err)
	// }

	if err := p.Render(music); err != nil {
		log.Fatal().Err(err)
	}
}
