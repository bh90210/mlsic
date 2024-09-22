package main

import (
	"flag"
	"fmt"
	"math/big"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/bh90210/mlsic"
	"github.com/bh90210/mlsic/markov"
	"github.com/bh90210/mlsic/markov/seed"
	"github.com/bh90210/mlsic/render"
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

	m := markov.Models{
		Poly: gomarkov.NewChain(2),
	}

	// Seed composition generation.
	poly := seedTrain()

	// Add the data to the model.
	m.AddPoly(poly)

	// Save seed model.
	err := m.Export(*modelsPath)
	if err != nil {
		log.Fatal().Err(err).Msg("exporting model")
	}

	noOfSpeakers := mlsic.TwoSpeakers

	// Generate the audio signal.
	speakers, err := markov.Deconstruct(poly, noOfSpeakers)
	if err != nil {
		log.Fatal().Err(err).Msg("deconstructing trains")
	}

	var music []mlsic.Audio
	music = append(music, speakers...)

	// Render audio to Port Audio.
	p, _ := render.NewPortAudio(render.WithChannels(noOfSpeakers))
	fmt.Println(*filesPath)

	if err := p.Render(music, "seed"); err != nil {
		log.Fatal().Err(err).Msg("rendering port audio")
	}

	// Render audio as .wav files.
	pp := render.Wav{
		Filepath: *filesPath,
	}

	if err := pp.Render(music, "seed"); err != nil {
		log.Fatal().Err(err).Msg("rendering wav files")
	}
}

// seedTrain .
func seedTrain() mlsic.Poly {
	log.Info().Msg("melody train")

	var poly mlsic.Poly

	voice1 := make(mlsic.Voice, 0)
	var prev int

	// 400.
	voice1[prev] = make(mlsic.Train)
	voice1[prev][0] = mlsic.Wagon{
		Sine: mlsic.Sine{
			Frequency: 400,
			Amplitude: 0.,
			Duration:  time.Duration(15 * time.Millisecond),
		},
		Panning: 0.,
	}

	for i := 0.02; i < 1.; i += 0.02 {
		prev += int(voice1[prev][0].Sine.Duration.Abs().Milliseconds()) *
			mlsic.SignalLengthMultiplier

		voice1[prev] = make(mlsic.Train)
		voice1[prev][0] = mlsic.Wagon{
			Sine: mlsic.Sine{
				Frequency: 400,
				Amplitude: i / 4,
				Duration:  time.Duration(15 * time.Millisecond),
			},
			Panning: i,
		}
	}

	for i := 0.98; i > 0.02; i -= 0.02 {
		prev += int(voice1[prev][0].Sine.Duration.Abs().Milliseconds()) *
			mlsic.SignalLengthMultiplier

		voice1[prev] = make(mlsic.Train)
		voice1[prev][0] = mlsic.Wagon{
			Sine: mlsic.Sine{
				Frequency: 400,
				Amplitude: i / 4,
				Duration:  time.Duration(30 * time.Millisecond),
			},
			Panning: i,
		}
	}

	h := &primeHarmonics{}
	seed.PartialsGeneration(h, voice1)

	voice2 := make(mlsic.Voice, 0)
	prev2 := prev

	// 400.
	voice2[prev2] = make(mlsic.Train)
	voice2[prev2][0] = mlsic.Wagon{
		Sine: mlsic.Sine{
			Frequency: 800,
			Amplitude: 0.,
			Duration:  time.Duration(30 * time.Millisecond),
		},
		Panning: 0.,
	}

	for i := 0.98; i > 0.02; i -= 0.02 {
		prev2 += int(voice2[prev2][0].Sine.Duration.Abs().Milliseconds()) *
			mlsic.SignalLengthMultiplier

		voice2[prev2] = make(mlsic.Train)
		voice2[prev2][0] = mlsic.Wagon{
			Sine: mlsic.Sine{
				Frequency: 500,
				Amplitude: i / 4,
				Duration:  time.Duration(15 * time.Millisecond),
			},
			Panning: i,
		}
	}

	for i := 0.02; i < 1.; i += 0.02 {
		prev2 += int(voice2[prev2][0].Sine.Duration.Abs().Milliseconds()) *
			mlsic.SignalLengthMultiplier

		voice2[prev2] = make(mlsic.Train)
		voice2[prev2][0] = mlsic.Wagon{
			Sine: mlsic.Sine{
				Frequency: 500,
				Amplitude: i / 4,
				Duration:  time.Duration(30 * time.Millisecond),
			},
			Panning: i,
		}
	}

	seed.PartialsGeneration(h, voice2)

	for i := 0.02; i < 1.; i += 0.02 {
		prev += int(voice1[prev][0].Sine.Duration.Abs().Milliseconds()) *
			mlsic.SignalLengthMultiplier

		voice1[prev] = make(mlsic.Train)
		voice1[prev][0] = mlsic.Wagon{
			Sine: mlsic.Sine{
				Frequency: 600,
				Amplitude: i / 4,
				Duration:  time.Duration(15 * time.Millisecond),
			},
			Panning: i,
		}
	}

	for i := 0.98; i > 0.02; i -= 0.02 {
		prev += int(voice1[prev][0].Sine.Duration.Abs().Milliseconds()) *
			mlsic.SignalLengthMultiplier

		voice1[prev] = make(mlsic.Train)
		voice1[prev][0] = mlsic.Wagon{
			Sine: mlsic.Sine{
				Frequency: 700,
				Amplitude: i / 4,
				Duration:  time.Duration(30 * time.Millisecond),
			},
			Panning: i,
		}
	}

	seed.PartialsGeneration(h, voice1)

	// Append voices to poly slice.
	// poly = append(poly, voice1, voice2)
	poly = append(poly, voice1)

	return poly
}

var _ mlsic.Harmonics = (*primeHarmonics)(nil)

type primeHarmonics struct {
	partials []mlsic.Partial
}

// Prime .
func (p *primeHarmonics) Partials() []mlsic.Partial {
	if len(p.partials) > 0 {
		return p.partials
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

	return p.partials
}

// func       amp          time
// ---------------------------------
// 400    0 -> 1 -> 0   2sec -> 4sec

// 800    0 -> 1 -> 0   2sec -> 4sec
// 200    0 -> 1 -> 0   2sec -> 4sec

// 1600    0 -> 1 -> 0   2sec -> 4sec
// 400    0 -> 1 -> 0   2sec -> 4sec
// 100    0 -> 1 -> 0   2sec -> 4sec

//
// Voice 1.
//

// voice1Trains := make(mlsic.Voice, 0)

// // Manual .
// voice1Trains[0] = make(mlsic.Train)

// voice1Trains[0][0] = mlsic.Wagon{
// 	Sine: mlsic.Sine{
// 		Frequency: 555,
// 		Amplitude: 0.49,
// 		Duration:  time.Duration(400 * time.Millisecond),
// 	},
// 	Panning: 0.,
// }

// var prev int

// prev += int(voice1Trains[prev][0].Sine.Duration.Abs().Milliseconds()) *
// 	mlsic.SignalLengthMultiplier

// voice1Trains[prev] = make(mlsic.Train)

// voice1Trains[prev][0] = mlsic.Wagon{
// 	Sine: mlsic.Sine{
// 		Frequency: 350,
// 		Amplitude: 0.4,
// 		Duration:  time.Duration(400 * time.Millisecond),
// 	},
// 	Panning: 0.25,
// }

// // Harmonics .
// var h mlsic.Harmonics
// h.PartialsGeneration(voice1Trains)

// // Append voice to poly slice.
// poly = append(poly, voice1Trains)

// //
// // Second voice.
// //

// voice2Trains := make(mlsic.Voice, 0)

// voice2Trains[0] = make(mlsic.Train)

// voice2Trains[0][0] = mlsic.Wagon{
// 	Sine: mlsic.Sine{
// 		Frequency: 1000,
// 		Amplitude: 0.29,
// 		Duration:  time.Duration(180 * time.Millisecond),
// 	},
// 	Panning: 0.50,
// }

// prev = 0

// prev += int(voice2Trains[prev][0].Sine.Duration.Abs().Milliseconds()) *
// 	mlsic.SignalLengthMultiplier

// voice2Trains[prev] = make(mlsic.Train)

// voice2Trains[prev][0] = mlsic.Wagon{
// 	Sine: mlsic.Sine{
// 		Frequency: 750,
// 		Amplitude: 0.2,
// 		Duration:  time.Duration(250 * time.Millisecond),
// 	},
// 	Panning: 0.75,
// }

// // Harmonics .
// h.PartialsGeneration(voice2Trains)
