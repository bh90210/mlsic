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

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
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

	// Move 1.
	var toneIndex int
	// toneIndex = upDown(80., .5, toneIndex, voice1, 0.1, 5000)

	// upDown(90., .6, toneIndex, voice1, 0.1, 1000)
	// toneIndex = upDown(80., .0, toneIndex, voice2, 0.1, 1000)

	// // Move 2.
	// toneIndex = upDown(80., .7, toneIndex, voice1, 0.0001, 5)

	// upDown(90., .8, toneIndex, voice1, 0.0001, 5)
	// toneIndex = upDown(80., .2, toneIndex, voice2, 0.0001, 5)

	// // Move 3.
	// move3(toneIndex, voice1, voice2, voice3, voice4)

	// Generate the partials.
	// voice1 = seed.Partials(voice1, primeMove1)
	// voice2 = seed.Partials(voice2, primeMove1)
	// voice3 = seed.Partials(voice3, primeMove1)
	// voice4 = seed.Partials(voice4, primeMove1)

	// Move 4.
	move4(toneIndex, voice1, voice2, voice3, voice4)

	var h seed.PrimeHarmonics
	h.PartialsGen(voice1)
	h.PartialsGen(voice2)
	h.PartialsGen(voice3)
	h.PartialsGen(voice4)

	// Append voices to poly slice.
	poly = append(poly, voice1, voice2, voice3, voice4)

	return poly
}

func upDown(freq float64, pan float64, tone int, voice markov.Voice, factor float64, duration int) int {
	for i := 0.; i < 1.; i += factor {
		tone += voice[tone].Fundamental.DurationInSamples()

		voice[tone] = markov.Tone{
			Fundamental: markov.Sine{
				Frequency: freq,
				Amplitude: i / 8,
				Duration:  time.Duration(duration) * time.Millisecond,
			},
			Panning: pan,
		}
	}

	for i := 1.; i > 0.; i -= 0.1 {
		tone += voice[tone].Fundamental.DurationInSamples()

		voice[tone] = markov.Tone{
			Fundamental: markov.Sine{
				Frequency: freq,
				Amplitude: i / 8,
				Duration:  time.Duration(5 * time.Millisecond),
			},
			Panning: pan,
		}
	}

	return tone
}

var primeMove1 = []mlsic.Partial{
	{
		Number:          2,
		AmplitudeFactor: .01,
		Start:           time.Duration(10 * time.Millisecond),
		Duration:        time.Duration(30 * time.Millisecond),
	},
	{
		Number:          3,
		AmplitudeFactor: .007,
		Start:           time.Duration(25 * time.Millisecond),
		Duration:        time.Duration(50 * time.Millisecond),
	},
	{
		// Number:          5,
		Number:          4,
		AmplitudeFactor: .06,
		Start:           time.Duration(3 * time.Millisecond),
		Duration:        time.Duration(100 * time.Millisecond),
	},
	{
		// Number:          7,
		Number:          5,
		AmplitudeFactor: .05,
		Start:           time.Duration(50 * time.Millisecond),
		Duration:        time.Duration(20 * time.Millisecond),
	},
	{
		// Number:          11,
		Number:          6,
		AmplitudeFactor: .05,
		Start:           time.Duration(30 * time.Millisecond),
		Duration:        time.Duration(30 * time.Millisecond),
	},
	{
		// Number:          13,
		Number:          7,
		AmplitudeFactor: .02,
		Start:           time.Duration(35 * time.Millisecond),
		Duration:        time.Duration(200 * time.Millisecond),
	},
	{
		// Number:          17,
		Number:          8,
		AmplitudeFactor: .015,
		Start:           time.Duration(100 * time.Millisecond),
		Duration:        time.Duration(100 * time.Millisecond),
	},
	{
		// Number:          19,
		Number:          9,
		AmplitudeFactor: .001,
		Start:           time.Duration(200 * time.Millisecond),
		Duration:        time.Duration(200 * time.Millisecond),
	},
	{
		// Number:          23,
		Number:          10,
		AmplitudeFactor: .005,
		Start:           time.Duration(2 * time.Millisecond),
		Duration:        time.Duration(2 * time.Millisecond),
	},
}

func move3(toneIndex int, voices ...markov.Voice) []markov.Voice {
	var freq float64 = 440.
	var pan float64 = .0
	var duration int = 5
	var factor1, factor2 float64 = .1, .1

	for i := 0; i < 10; i++ {
		for _, voice := range voices {
			move3UpDown(freq, pan, toneIndex, voice, factor1, factor2, duration)

			freq += 100.

			pan += .1
			if pan > 1. {
				pan = .0
			}

			duration++
		}
	}

	for _, v := range voices {
		if v.LengthInSamples() > toneIndex {
			toneIndex = v.LengthInSamples()
		}
	}

	freq = 4440.
	duration = 5
	factor1, factor2 = .05, .01

	for i := 0; i < 10; i++ {
		for voiceIndex, voice := range voices {
			newIndex := move3UpDown(freq, pan, toneIndex, voice, factor1, factor2, duration)

			freq -= 10.

			pan += .01
			if pan > 1. {
				pan = .0
			}

			duration++

			if len(voices) == voiceIndex {
				toneIndex = newIndex
			}
		}
	}

	for _, v := range voices {
		if v.LengthInSamples() > toneIndex {
			toneIndex = v.LengthInSamples()
		}
	}

	duration = 10
	factor1, factor2 = .5, .1

	for i := 0; i < 30; i++ {
		for voiceIndex, voice := range voices {
			newIndex := move3UpDown(freq, pan, toneIndex, voice, factor1, factor2, duration)

			freq += 5.

			pan += .01
			if pan > 1. {
				pan = .0
			}

			// duration++

			if len(voices)-1 == voiceIndex {
				toneIndex = newIndex
			}
		}
	}

	toneIndex += 2000
	duration = 10
	freq = 100
	factor1, factor2 = .1, .01

	for i := 0; i < 3; i++ {
		for voiceIndex, voice := range voices {
			newIndex := move3UpDown(freq, pan, toneIndex, voice, factor1, factor2, duration)

			freq += 50.

			pan += .05
			if pan > 1. {
				pan = .0
			}

			// duration++

			if len(voices)-1 == voiceIndex {
				toneIndex = newIndex
			}
		}
	}

	duration = 10
	factor1, factor2 = .5, .1

	for i := 0; i < 20; i++ {
		for voiceIndex, voice := range voices {
			newIndex := move3UpDown(freq, pan, toneIndex, voice, factor1, factor2, duration)

			freq += 5.

			pan += .01
			if pan > 1. {
				pan = .0
			}

			// duration++

			if len(voices)-1 == voiceIndex {
				toneIndex = newIndex
			}
		}
	}

	return voices
}

func move3UpDown(freq float64, pan float64, tone int, voice markov.Voice, factor1, factor2 float64, duration int) int {
	for i := 0.; i < 1.; i += factor1 {
		tone += voice[tone].Fundamental.DurationInSamples()

		voice[tone] = markov.Tone{
			Fundamental: markov.Sine{
				Frequency: freq,
				Amplitude: i / 8,
				Duration:  time.Duration(duration) * time.Millisecond,
			},
			Panning: pan,
		}
	}

	for i := 1.; i > 0.; i -= factor2 {
		tone += voice[tone].Fundamental.DurationInSamples()

		voice[tone] = markov.Tone{
			Fundamental: markov.Sine{
				Frequency: freq,
				Amplitude: i / 8,
				Duration:  time.Duration(5 * time.Millisecond),
			},
			Panning: pan,
		}
	}

	return tone
}

func move4(toneIndex int, voices ...markov.Voice) []markov.Voice {
	for _, voice := range voices {
		toneIndex += voice[toneIndex].Fundamental.DurationInSamples()

		voice[toneIndex] = markov.Tone{
			Fundamental: markov.Sine{
				Frequency: 1000,
				Amplitude: .2 / 8,
				Duration:  time.Duration(1000) * time.Millisecond,
			},
			Panning: 0.5,
		}
	}

	for voiceIndex, voice := range voices {
		var freq, pan float64

		switch voiceIndex {
		case 0, 1:
			freq = 1000
			pan = 0.5

		case 2:
			freq = 900
			pan = 0.

		case 3:
			freq = 1100
			pan = 1.

		}

		toneIndex += voice[toneIndex].Fundamental.DurationInSamples()

		voice[toneIndex] = markov.Tone{
			Fundamental: markov.Sine{
				Frequency: freq,
				Amplitude: .2 / 8,
				Duration:  time.Duration(1000) * time.Millisecond,
			},
			Panning: pan,
		}
	}

	for voiceIndex, voice := range voices {
		var freq, pan float64

		switch voiceIndex {
		case 0:
			freq = 900
			pan = 0.4

		case 1:
			freq = 1100
			pan = 0.6

		case 2:
			freq = 800
			pan = 0.

		case 3:
			freq = 1200
			pan = 1.

		}

		toneIndex += voice[toneIndex].Fundamental.DurationInSamples()

		voice[toneIndex] = markov.Tone{
			Fundamental: markov.Sine{
				Frequency: freq,
				Amplitude: .2 / 8,
				Duration:  time.Duration(1000) * time.Millisecond,
			},
			Panning: pan,
		}
	}

	for voiceIndex, voice := range voices {
		var freq, pan float64

		switch voiceIndex {
		case 0:
			freq = 950
			pan = 0.4

		case 1:
			freq = 1050
			pan = 0.6

		case 2:
			freq = 700
			pan = 0.

		case 3:
			freq = 1300
			pan = 1.

		}

		toneIndex += voice[toneIndex].Fundamental.DurationInSamples()

		voice[toneIndex] = markov.Tone{
			Fundamental: markov.Sine{
				Frequency: freq,
				Amplitude: .2 / 8,
				Duration:  time.Duration(1000) * time.Millisecond,
			},
			Panning: pan,
		}
	}

	return voices
}
