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

	m := markov.Model{
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
func polySeed() []markov.Voice {
	log.Info().Msg("melody train")

	var poly []markov.Voice

	voice1 := make(markov.Voice)
	voice2 := make(markov.Voice)
	voice3 := make(markov.Voice)
	voice4 := make(markov.Voice)

	var toneIndex int
	toneIndex = upDown(400., .5, toneIndex, voice1)

	upDown(500., .6, toneIndex, voice1)
	toneIndex = upDown(300., .4, toneIndex, voice2)

	upDown(600., .7, toneIndex, voice1)
	upDown(200., .3, toneIndex, voice2)
	toneIndex = upDown(400., .5, toneIndex, voice3)

	upDown(650., .65, toneIndex, voice1)
	upDown(150., .2, toneIndex, voice2)
	upDown(350., .35, toneIndex, voice3)
	toneIndex = upDown(450., .4, toneIndex, voice4)

	upDown(700., .65, toneIndex, voice1)
	upDown(100., .2, toneIndex, voice2)
	upDown(300., .35, toneIndex, voice3)
	toneIndex = upDown(500., .4, toneIndex, voice4)

	upDown(690., .65, toneIndex, voice1)
	upDown(110., .2, toneIndex, voice2)
	upDown(350., .35, toneIndex, voice3)
	toneIndex = upDown(550., .4, toneIndex, voice4)

	// Append voices to poly slice.
	poly = append(poly, voice1, voice2, voice3, voice4)

	// Generate the partials.
	h := seed.PrimeHarmonics{}
	poly = h.Partials(poly)

	return poly
}

func upDown(freq float64, pan float64, tone int, voice markov.Voice) int {
	for i := 0.; i < 1.; i += 0.1 {
		tone += voice[tone].Fundamental.DurationInSamples()

		voice[tone] = markov.Tone{
			Fundamental: markov.Sine{
				Frequency: freq,
				Amplitude: i / 4,
				Duration:  time.Duration(5 * time.Millisecond),
			},
			Panning: pan,
		}
	}

	for i := 1.; i > 0.; i -= 0.1 {
		tone += voice[tone].Fundamental.DurationInSamples()

		voice[tone] = markov.Tone{
			Fundamental: markov.Sine{
				Frequency: freq,
				Amplitude: i / 4,
				Duration:  time.Duration(5 * time.Millisecond),
			},
			Panning: pan,
		}
	}

	return tone
}
