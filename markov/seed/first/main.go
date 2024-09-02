package main

import (
	"flag"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/bh90210/mlsic"
	"github.com/bh90210/mlsic/markov"
	"github.com/bh90210/mlsic/render"
	"github.com/go-audio/generator"
	"github.com/mb-14/gomarkov"
)

func main() {
	debug := flag.Bool("debug", false, "sets log level to debug")
	filesPath := flag.String("files", "", "sets the directory audio files will be saved")
	modelsPath := flag.String("models", "", "sets the directory model files will be saved")

	flag.Parse()

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	m := markov.Train{
		Freq: gomarkov.NewChain(1),
		Amp:  gomarkov.NewChain(1),
		Dur:  gomarkov.NewChain(1),
	}

	var from int

	// Composition.
	var train []mlsic.Sine
	for i := 0; i < 50; i++ {
		train = append(train, mlsic.Sine{
			Frequency: (440 + (float64(i) * 2)) / float64(i),
			Amplitude: 0.07,
			Duration:  time.Duration((100. + float64(i))) * time.Millisecond,
		})
	}

	m.Add(train)
	from = len(train)

	for i := 0; i < 69; i++ {
		train = append(train, mlsic.Sine{
			Frequency: 440. / float64(i),
			Amplitude: 0.003 * float64(i),
			Duration:  time.Duration(205.-float64(i)) * time.Millisecond,
		},
			mlsic.Sine{
				Frequency: 440.,
				Amplitude: 0.,
				Duration:  time.Duration(5.+float64(i)) * time.Millisecond,
			})
	}

	m.Add(train[from:])
	from = len(train)

	for i := 0; i < 69; i++ {
		train = append(train, mlsic.Sine{
			Frequency: (440 + (float64(i) * 2)) / float64(i),
			Amplitude: 0.01 * float64(i),
			Duration:  time.Duration(100.) * time.Millisecond,
		})
	}

	m.Add(train[from:])
	from = len(train)

	for i := 0; i < 20; i++ {
		train = append(train, mlsic.Sine{
			Frequency: 100. * float64(i),
			Amplitude: 0.2,
			Duration:  time.Duration(10.) * time.Millisecond,
		})
	}

	m.Add(train[from:])
	from = len(train)

	for i := 0; i < 20; i++ {
		train = append(train, mlsic.Sine{
			Frequency: 200. * float64(i),
			Amplitude: 0.2,
			Duration:  time.Duration(10.) * time.Millisecond,
		})
	}

	m.Add(train[from:])
	from = len(train)

	for i := 0; i < 20; i++ {
		train = append(train, mlsic.Sine{
			Frequency: 300. * float64(i),
			Amplitude: 0.2,
			Duration:  time.Duration(10.) * time.Millisecond,
		})
	}

	m.Add(train[from:])
	from = len(train)

	for i := 0; i < 20; i++ {
		train = append(train, mlsic.Sine{
			Frequency: 400. * float64(i),
			Amplitude: 0.2,
			Duration:  time.Duration(10.) * time.Millisecond,
		})
	}

	m.Add(train[from:])
	from = len(train)

	for i := 0; i < 10; i++ {
		train = append(train, mlsic.Sine{
			Frequency: 400. * float64(i),
			Amplitude: 0.1,
			Duration:  time.Duration(100.) * time.Millisecond,
		})
	}

	m.Add(train[from:])
	from = len(train)

	for i := 0; i < 10; i++ {
		train = append(train, mlsic.Sine{
			Frequency: 750. * float64(i),
			Amplitude: 0.1,
			Duration:  time.Duration(100.) * time.Millisecond,
		})
	}

	m.Add(train[from:])
	from = len(train)

	for i := 0; i < 15; i++ {
		train = append(train, mlsic.Sine{
			Frequency: 250. * float64(i),
			Amplitude: 0.1,
			Duration:  time.Duration(100.) * time.Millisecond,
		})
	}

	m.Add(train[from:])
	from = len(train)

	for i := 0; i < 10; i++ {
		train = append(train, mlsic.Sine{
			Frequency: 350. * float64(i),
			Amplitude: 0.1,
			Duration:  time.Duration(70.) * time.Millisecond,
		})
	}

	m.Add(train[from:])
	from = len(train)

	for i := 0; i < 5; i++ {
		train = append(train, mlsic.Sine{
			Frequency: 350. * float64(i),
			Amplitude: 0.,
			Duration:  time.Duration(70.) * time.Millisecond,
		})
	}

	m.Add(train[from:])
	from = len(train)

	for i := 0; i < 10; i++ {
		train = append(train, mlsic.Sine{
			Frequency: 400. * float64(i),
			Amplitude: 0.1,
			Duration:  time.Duration(90.) * time.Millisecond,
		})
	}

	m.Add(train[from:])
	from = len(train)

	for i := 0; i < 10; i++ {
		train = append(train, mlsic.Sine{
			Frequency: 750. * float64(i),
			Amplitude: 0.1,
			Duration:  time.Duration(110.) * time.Millisecond,
		})
	}

	m.Add(train[from:])
	from = len(train)

	for i := 0; i < 10; i++ {
		train = append(train, mlsic.Sine{
			Frequency: 250. * float64(i),
			Amplitude: 0.1,
			Duration:  time.Duration(120.) * time.Millisecond,
		})
	}

	m.Add(train[from:])
	from = len(train)

	for i := 0; i < 10; i++ {
		train = append(train, mlsic.Sine{
			Frequency: 200. * float64(i),
			Amplitude: 0.1,
			Duration:  time.Duration(120.) * time.Millisecond,
		})
	}

	m.Add(train[from:])
	from = len(train)

	for i := 0; i < 10; i++ {
		train = append(train, mlsic.Sine{
			Frequency: 150. * float64(i),
			Amplitude: 0.1,
			Duration:  time.Duration(120.) * time.Millisecond,
		})
	}

	m.Add(train[from:])
	from = len(train)

	for i := 0; i < 10; i++ {
		train = append(train, mlsic.Sine{
			Frequency: 100. * float64(i),
			Amplitude: 0.1,
			Duration:  time.Duration(120.) * time.Millisecond,
		})
	}

	m.Add(train[from:])
	from = len(train)

	for i := 0; i < 10; i++ {
		train = append(train, mlsic.Sine{
			Frequency: 50. * float64(i),
			Amplitude: 0.1,
			Duration:  time.Duration(120.) * time.Millisecond,
		})
	}

	m.Add(train[from:])
	from = len(train)

	for i := 0; i < 10; i++ {
		train = append(train, mlsic.Sine{
			Frequency: 0. * float64(i),
			Amplitude: 0.1,
			Duration:  time.Duration(120.) * time.Millisecond,
		})
	}

	m.Add(train[from:])
	from = len(train)

	for i := 0; i < 10; i++ {
		train = append(train, mlsic.Sine{
			Frequency: 1. * float64(i),
			Amplitude: 0.1,
			Duration:  time.Duration(120.) * time.Millisecond,
		})
	}

	m.Add(train[from:])
	from = len(train)

	for i := 0; i < 10; i++ {
		train = append(train, mlsic.Sine{
			Frequency: 1. * float64(i),
			Amplitude: 0.,
			Duration:  time.Duration(120.) * time.Millisecond,
		})
	}

	m.Add(train[from:])
	from = len(train)

	// Weird break.

	for i := 1; i < 5; i++ {
		train = append(train, mlsic.Sine{
			Frequency: 1. * float64(i),
			Amplitude: 0.1,
			Duration:  time.Duration(5000.*i) * time.Millisecond,
		})
	}

	m.Add(train[from:])
	from = len(train)

	for i := 0; i < 10; i++ {
		train = append(train, mlsic.Sine{
			Frequency: 1. * float64(i),
			Amplitude: 0.,
			Duration:  time.Duration(120.) * time.Millisecond,
		})
	}

	m.Add(train[from:])
	from = len(train)

	for i := 0; i < 10; i++ {
		train = append(train,
			mlsic.Sine{
				Frequency: 1. * float64(i),
				Amplitude: 0.1,
				Duration:  time.Duration(108.) * time.Millisecond,
			},
			mlsic.Sine{
				Frequency: 250. * float64(i),
				Amplitude: 0.1,
				Duration:  time.Duration(12.) * time.Millisecond,
			})
	}

	m.Add(train[from:])
	from = len(train)

	for i := 0; i < 10; i++ {
		train = append(train,
			mlsic.Sine{
				Frequency: 1. * float64(i),
				Amplitude: 0.1,
				Duration:  time.Duration(100.) * time.Millisecond,
			},
			mlsic.Sine{
				Frequency: 50. * float64(i),
				Amplitude: 0.1,
				Duration:  time.Duration(50.) * time.Millisecond,
			})
	}

	m.Add(train[from:])
	from = len(train)

	// Save base model.
	err := m.Export(*modelsPath)
	if err != nil {
		log.Fatal().Err(err).Msg("exporting models")
	}

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
	for _, v := range train {
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

	// Render audio as .wav files.
	p := render.Wav{
		Filepath: *filesPath,
	}

	if err := p.Render(music, "seed"); err != nil {
		log.Fatal().Err(err).Msg("rendering")
	}
}
