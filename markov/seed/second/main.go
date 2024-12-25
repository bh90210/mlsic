package main

import (
	"flag"
	"os"

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

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	m := &markov.Model{
		Freq: gomarkov.NewChain(1),
		Dur:  gomarkov.NewChain(1),
		Amp:  gomarkov.NewChain(1),
		Pan:  gomarkov.NewChain(1),
	}

	// Seed composition generation.
	poly := polySeed()

	m.Add2(poly)

	// Save seed model.
	err := m.Export(*modelsPath)
	if err != nil {
		log.Fatal().Err(err).Msg("exporting model")
	}

	// Generate the audio signal for each speakers (mlsic.TwoSpeakers.)
	speakersSignal, err := markov.Deconstruct(poly, mlsic.TwoSpeakers)
	if err != nil {
		log.Fatal().Err(err).Msg("deconstructing trains")
	}

	var music []mlsic.Audio
	music = append(music, speakersSignal...)

	// Render audio to Port Audio.
	// portAudio, _ := render.NewPortAudio(render.WithChannels(mlsic.TwoSpeakers))

	// if err := portAudio.Render(music, "seed"); err != nil {
	// 	log.Fatal().Err(err).Msg("rendering port audio")
	// }

	// fmt.Println(*filesPath)
	// Render audio as .wav files.
	wav := render.Wav{
		Filepath: *filesPath,
	}

	if err := wav.Render(music, "seed"); err != nil {
		log.Fatal().Err(err).Msg("rendering wav files")
	}
}
