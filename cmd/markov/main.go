package main

import (
	"flag"
	"os"

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

	s := markov.Song{
		NGenerations:  *ngenerations,
		FilePath:      *filesPath,
		ModelsPath:    *modelsPath,
		SeedModelPath: *seedModelPath,
	}

	s.NGen()
}
