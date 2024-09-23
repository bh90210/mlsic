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
	poly := polySeed()

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

// polySeed .
func polySeed() mlsic.Poly {
	log.Info().Msg("melody train")

	var poly mlsic.Poly

	voice1 := make(mlsic.Voice, 0)
	var prev int

	// 400.
	for i := 0.; i < 1.; i += 0.01 {
		prev += voice1[prev][0].Sine.DurationInSamples()

		voice1[prev] = make(mlsic.Train)
		voice1[prev][0] = mlsic.Wagon{
			Sine: mlsic.Sine{
				Frequency: 440.,
				Amplitude: i / 4,
				Duration:  time.Duration(30 * time.Millisecond),
			},
			Panning: i,
		}
	}

	for i := 1.; i > 0.; i -= 0.01 {
		prev += voice1[prev][0].Sine.DurationInSamples()

		voice1[prev] = make(mlsic.Train)
		voice1[prev][0] = mlsic.Wagon{
			Sine: mlsic.Sine{
				Frequency: 440.,
				Amplitude: i / 4,
				Duration:  time.Duration(30 * time.Millisecond),
			},
			Panning: i,
		}
	}

	h := &primeHarmonics{}
	seed.PartialsGeneration(h, voice1)

	// Append voices to poly slice.
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
