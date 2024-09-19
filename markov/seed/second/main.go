package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/bh90210/mlsic"
	"github.com/bh90210/mlsic/markov"
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
	poly := MelodyTrain()

	// Add the data to the model.
	m.AddPoly(poly)

	// Save seed model.
	err := m.Export(*modelsPath)
	if err != nil {
		log.Fatal().Err(err).Msg("exporting model")
	}

	channels, err := markov.DeconstructTrains(poly, 2)
	if err != nil {
		log.Fatal().Err(err).Msg("deconstructing trains")
	}

	var music []mlsic.Audio
	music = append(music, channels...)

	p, _ := render.NewPortAudio()
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

// MelodyTrain .
func MelodyTrain() markov.Poly {
	log.Info().Msg("melody train")

	voice1Trains := make(markov.Voice, 0)

	// Manual .
	voice1Trains[0] = make(markov.Train)

	voice1Trains[0][0] = markov.Wagon{
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

	voice1Trains[prev][0] = markov.Wagon{
		Sine: mlsic.Sine{
			Frequency: 350,
			Amplitude: 0.4,
			Duration:  time.Duration(400 * time.Millisecond),
		},
		Panning: 0.25,
	}

	// Harmonics .
	var h markov.Harmonics
	h.PartialsGeneration(voice1Trains)

	var poly markov.Poly
	poly = append(poly, voice1Trains)

	//
	// Second voice.
	//

	voice2Trains := make(markov.Voice, 0)

	// Manual .
	voice2Trains[0] = make(markov.Train)

	voice2Trains[0][0] = markov.Wagon{
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

	voice2Trains[prev][0] = markov.Wagon{
		Sine: mlsic.Sine{
			Frequency: 750,
			Amplitude: 0.2,
			Duration:  time.Duration(250 * time.Millisecond),
		},
		Panning: 0.75,
	}

	// Harmonics .
	h.PartialsGeneration(voice2Trains)

	poly = append(poly, voice2Trains)

	return poly
}
