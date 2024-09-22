package main

import (
	"flag"
	"os"

	"github.com/bh90210/mlsic"
	"github.com/bh90210/mlsic/markov"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	debug := flag.Bool("debug", false, "sets log level to debug")
	ngenerations := flag.Int("ngen", 2, "sets log level to debug")
	filesPath := flag.String("files", "", "sets the directory audio files will be saved")
	modelsPath := flag.String("models", "", "sets the directory model files will be saved")
	seedModelPath := flag.String("seed", "", "sets the directory of seed model to use")

	flag.Parse()

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	// Init a markov song.
	s := markov.Song{
		NGenerations:  *ngenerations,
		FilePath:      *filesPath,
		ModelsPath:    *modelsPath,
		SeedModelPath: *seedModelPath,
		Harmonics:     &naive{},
	}

	s.NGen()
}

var _ mlsic.Harmonics = (*naive)(nil)

// naive .
type naive struct {
	partials []mlsic.Partial
}

// Partials .
func (n *naive) Partials() []mlsic.Partial {
	if len(n.partials) > 0 {
		return n.partials
	}

	// Harmonics.
	for i := 2; i < 180; i++ {
		n.partials = append(n.partials, mlsic.Partial{
			Number:          i,
			AmplitudeFactor: float64(i) * 0.01,
		})
	}

	return n.partials
}
