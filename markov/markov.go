package markov

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/bh90210/mlsic"
	"github.com/bh90210/mlsic/render"
	"github.com/go-audio/generator"
	"github.com/mb-14/gomarkov"
)

// Sine .
type Sine struct {
	Frequency float64
	Amplitude float64
	Duration  time.Duration
}

// Markov .
type Markov struct {
	Freq *gomarkov.Chain
	Amp  *gomarkov.Chain
	Dur  *gomarkov.Chain
}

// Add .
func (m *Markov) Add(train []Sine) {
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
func (m *Markov) Export(path string) error {
	freq, err := m.Freq.MarshalJSON()
	if err != nil {
		return err
	}

	err = os.WriteFile(path+"/freq.json", freq, 0644)
	if err != nil {
		return err
	}

	amp, err := m.Amp.MarshalJSON()
	if err != nil {
		return err
	}

	err = os.WriteFile(path+"/amp.json", amp, 0644)
	if err != nil {
		return err
	}

	dur, err := m.Dur.MarshalJSON()
	if err != nil {
		return err
	}

	err = os.WriteFile(path+"/dur.json", dur, 0644)
	if err != nil {
		return err
	}

	return nil
}

// Train .
type Train struct {
	Train []Sine

	Freqs [][]float64
	Amps  [][]float64
	Durs  [][]float64
}

// Field .
type Field int

const (
	// Freq .
	Freq Field = iota
	// Amp .
	Amp
	// Dur .
	Dur
)

// AddSine .
func (t *Train) AddSine(newTrain []Sine) {
	t.Train = append(t.Train, newTrain...)
}

// Add .
func (t *Train) Add(field Field, values []float64) {
	switch field {
	case Freq:
		t.Freqs = append(t.Freqs, values)
	case Amp:
		t.Amps = append(t.Amps, values)
	case Dur:
		t.Durs = append(t.Durs, values)
	}

	log.Debug().Any("length", len(values)).Any("field", field).Msg("add train")
}

// Generate .
func (t *Train) Generate(filepath string) {
	// Harmonics.
	partials := make(map[int]float64)
	for i := 2; i < 180; i++ {
		partials[i] = float64(i) * 0.01
		if partials[i] > 1. {
			partials[i] -= 1.
		}
	}

	// Left channel.
	var left []float64
	// Right channel.
	var right []float64
	for _, v := range t.Train {
		osc := generator.NewOsc(generator.WaveSine, v.Frequency, 44100)
		osc.Amplitude = v.Amplitude
		// osc.SetAttackInMs(10)

		signal := osc.Signal(44 * int(v.Duration.Abs().Milliseconds()))

		for partial, amplitude := range partials {
			if v.Frequency*float64(partial) > 18000 {
				continue
			}

			osc = generator.NewOsc(generator.WaveSine, v.Frequency*float64(partial), 44100)

			amplitude = amplitude + ((amplitude - partials[partial]) / 2)
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
	// 	log.Fatal(err)
	// }

	if err := p.Render(music); err != nil {
		log.Fatal().Err(err)
	}
}

type model struct {
	Int      int `json:"int"`
	SpoolMap any `json:"spool_map"`
	FreqMat  any `json:"freq_mat"`
}

// ErrYo .
var ErrYo = errors.New("test")

// Song .
type Song struct {
	NGenerations  int
	FilePath      string
	ModelsPath    string
	SeedModelPath string
}

// NGen .
func (s *Song) NGen() {
	log.Info().Msg("NGen")

	for i := 0; i < s.NGenerations; i++ {
		log.Info().Int("generation", i).Msg("NGen")

		var freqModel string
		var ampModel string
		var durModel string

		index := strconv.Itoa(i)

		if i == 0 {
			freqModel = filepath.Join(s.SeedModelPath, "freq.json")
			ampModel = filepath.Join(s.SeedModelPath, "amp.json")
			durModel = filepath.Join(s.SeedModelPath, "dur.json")
		} else {
			previousModelsIndex := strconv.Itoa(i - 1)

			freqModel = filepath.Join(s.ModelsPath, "gen"+previousModelsIndex, "freq.json")
			ampModel = filepath.Join(s.ModelsPath, "gen"+previousModelsIndex, "amp.json")
			durModel = filepath.Join(s.ModelsPath, "gen"+previousModelsIndex, "dur.json")
		}

		log.Info().Int("generation", i).Msg("reading files")

		freq, err := os.ReadFile(freqModel)
		if err != nil {
			log.Fatal().Err(err).Msg("reading freq")
		}

		amp, err := os.ReadFile(ampModel)
		if err != nil {
			log.Fatal().Err(err).Msg("reading amp")
		}

		dur, err := os.ReadFile(durModel)
		if err != nil {
			log.Fatal().Err(err).Msg("reading dur")
		}

		log.Info().Int("generation", i).Msg("creating chains")

		f := gomarkov.NewChain(1)
		f.UnmarshalJSON(freq)

		a := gomarkov.NewChain(1)
		a.UnmarshalJSON(amp)

		d := gomarkov.NewChain(1)
		d.UnmarshalJSON(dur)

		var t Train

		// Freq.
		var frequencies model
		err = json.Unmarshal(freq, &frequencies)
		if err != nil {
			log.Fatal().Err(err).Msg("unmarsh freq")
		}

		log.Info().Int("generation", i).Int("field", int(Freq)).Msg("entering loop")

		err = LoopMapped(frequencies.SpoolMap, &t, f, Freq)
		if err != nil {
			log.Fatal().Err(err).Msg("freq loop")
		}

		// Amp.
		var amplitudes model
		err = json.Unmarshal(amp, &amplitudes)
		if err != nil {
			log.Fatal().Err(err).Msg("unmarsh amp")
		}

		log.Info().Int("generation", i).Int("field", int(Amp)).Msg("entering loop")

		err = LoopMapped(amplitudes.SpoolMap, &t, a, Amp)
		if err != nil {
			if errors.Is(err, ErrYo) {
				fmt.Println("mama")
			}
			log.Fatal().Err(err).Msg("amp loop")
		}

		// Dur.
		var durations model
		err = json.Unmarshal(dur, &durations)
		if err != nil {
			log.Fatal().Err(err).Msg("unmarsh dur")
		}

		log.Info().Int("generation", i).Int("field", int(Dur)).Msg("entering loop")

		err = LoopMapped(durations.SpoolMap, &t, d, Dur)
		if err != nil {
			log.Fatal().Err(err).Msg("dur loop")
		}

		log.Info().Int("generation", i).Msg("creating sines train")

		// Create sines train.
		for i, freqs := range t.Freqs {
			var outOfBoundsAmp bool
			var outOfBoundsDur bool

			if len(t.Amps)-1 < i {
				outOfBoundsAmp = true
			}

			if len(t.Durs)-1 < i {
				outOfBoundsDur = true
			}

			for o, freq := range freqs {
				// TODO: 0 is arbitrary, fix it.
				amp := 0.
				if !outOfBoundsAmp && !(len(t.Amps[i])-1 < o) {
					amp = t.Amps[i][o]
				}

				// TODO: 100 is arbitrary, fix it.
				dur := 100.
				if !outOfBoundsDur && !(len(t.Durs[i])-1 < o) {
					dur = t.Durs[i][o]
				}

				t.Train = append(t.Train, Sine{
					Frequency: freq,
					Amplitude: amp,
					Duration:  time.Duration(dur) * time.Millisecond,
				})
			}
		}

		log.Info().Int("generation", i).Msg("audio files gen")

		filePath := filepath.Join(s.FilePath, "gen"+index)

		err = os.MkdirAll(filePath, 0755)
		if err != nil {
			log.Fatal().Err(err).Msg("creating audio directory")
		}

		t.Generate(filePath)

		log.Info().Int("generation", i).Msg("export models")

		// Save model.
		m := Markov{
			Freq: gomarkov.NewChain(1),
			Amp:  gomarkov.NewChain(1),
			Dur:  gomarkov.NewChain(1),
		}

		m.Add(t.Train)

		modelsPath := filepath.Join(s.ModelsPath, "gen"+index)

		err = os.MkdirAll(modelsPath, 0755)
		if err != nil {
			log.Fatal().Err(err).Msg("creating models directory")
		}

		err = m.Export(modelsPath)
		if err != nil {
			log.Fatal().Err(err).Msg("exporting seed")
		}
	}
}

// LoopMapped .
func LoopMapped(spoolMap any, t *Train, chain *gomarkov.Chain, field Field) error {
	mapped := spoolMap.(map[string]interface{})

	for key := range mapped {
		if key == "$" || key == "^" || key == "+Inf" {
			continue
		}

		var temporaryTrain []float64

		starting := []string{key}
		// TODO: 100 is arbitrary, fix it.
		for i := 0; i < 100; i++ {
			generated, err := chain.Generate(starting)
			if err != nil {
				return fmt.Errorf("generated %w %w %v %v", err, ErrYo, starting, i)
			}

			if generated == "$" {
				break
			}

			flo, err := strconv.ParseFloat(generated, 64)
			if err != nil {
				return fmt.Errorf("parse float %w", err)
			}

			temporaryTrain = append(temporaryTrain, flo)

			starting = []string{generated}
		}

		t.Add(field, temporaryTrain)
	}

	return nil
}
