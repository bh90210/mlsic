package main

import (
	"flag"
	"fmt"
	"os"

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
		Freq: gomarkov.NewChain(1),
		Amp:  gomarkov.NewChain(1),
		Dur:  gomarkov.NewChain(1),
	}

	// Seed composition generation.
	poly := seed.MelodyTrain()

	// Save seed model.
	err := m.Export(*modelsPath)
	if err != nil {
		log.Fatal().Err(err).Msg("exporting models")
	}

	channels, err := seed.DeconstructTrains(poly, 2)
	if err != nil {
		log.Fatal().Err(err).Msg("deconstructing trains")
	}

	var music []mlsic.Audio
	music = append(music, channels...)

	p, _ := render.NewPortAudio()
	fmt.Println(*filesPath)

	if err := p.Render(music, "seed"); err != nil {
		log.Fatal().Err(err).Msg("rendering")
	}

	// Render audio as .wav files.
	pp := render.Wav{
		Filepath: *filesPath,
	}

	if err := pp.Render(music, "seed"); err != nil {
		log.Fatal().Err(err).Msg("rendering")
	}
}
