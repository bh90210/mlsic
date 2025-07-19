package main

import (
	"flag"
	"os"
	"time"

	"github.com/bh90210/mlsic"
	"github.com/bh90210/mlsic/render"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	debug := flag.Bool("debug", false, "sets log level to debug")
	// filesPath := flag.String("files", "", "sets the directory audio files will be saved")

	flag.Parse()

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	// Seed composition generation.
	poly := polySeed()

	// Generate the audio signal for each speakers (mlsic.TwoSpeakers.)
	speakersSignal, err := mlsic.Signal(poly, mlsic.TwoSpeakers)
	if err != nil {
		log.Fatal().Err(err).Msg("deconstructing trains")
	}

	var music []mlsic.Audio
	music = append(music, speakersSignal...)

	// Render audio to Port Audio.
	portAudio, _ := render.NewPortAudio(render.WithChannels(mlsic.TwoSpeakers))

	if err := portAudio.Render(music, "seed"); err != nil {
		log.Fatal().Err(err).Msg("rendering port audio")
	}

	// fmt.Println(*filesPath)
	// Render audio as .wav files.
	// wav := render.Wav{
	// 	Filepath: *filesPath,
	// }

	// if err := wav.Render(music, "seed"); err != nil {
	// 	log.Fatal().Err(err).Msg("rendering wav files")
	// }
}

// polySeed .
func polySeed() []mlsic.Voice {
	log.Info().Msg("melody train")

	var poly []mlsic.Voice

	voice1 := make(mlsic.Voice)
	voice2 := make(mlsic.Voice)
	voice3 := make(mlsic.Voice)
	voice4 := make(mlsic.Voice)

	// Move 1.
	var toneIndex int

	toneIndex = upDown(280., .5, toneIndex, voice1, 0.1, 50)

	upDown(190., .6, toneIndex, voice1, 0.1, 50)

	toneIndex = upDown(380., .0, toneIndex, voice2, 0.1, 50)

	// Move 2.
	toneIndex = upDown(480., .7, toneIndex, voice1, 0.01, 20)

	upDown(90., .8, toneIndex, voice1, 0.01, 20)

	toneIndex = upDown(580., .2, toneIndex, voice2, 0.01, 20)

	// Move 3.
	move3(toneIndex, voice1, voice2, voice3, voice4)

	// Move 4.
	move4(voice1.LengthInSamples(), voice1, voice2, voice3, voice4)

	// Generate the partials.
	voice1 = mlsic.Partials(voice1, mlsic.PrimeMove1)
	voice2 = mlsic.Partials(voice2, mlsic.PrimeMove1)
	voice3 = mlsic.Partials(voice3, mlsic.PrimeMove1)
	voice4 = mlsic.Partials(voice4, mlsic.PrimeMove1)

	// Append voices to poly slice.
	poly = append(poly, voice1, voice2, voice3, voice4)

	return poly
}

func injectStart(index int, voice mlsic.Voice) int {
	voice[index] = mlsic.Tone{
		Fundamental: mlsic.Sine{
			Frequency: 180.000000,
			Amplitude: 0.100000,
			Duration:  time.Duration(50) * time.Millisecond,
		},
		Panning: 0.100000,
	}

	index += voice[index].Fundamental.DurationInSamples()

	return index
}

func upDown(freq float64, pan float64, toneIndex int, voice mlsic.Voice, factor float64, duration int) int {
	for i := 0.; i < 1.; i += factor {
		toneIndex += voice[toneIndex].Fundamental.DurationInSamples()

		voice[toneIndex] = mlsic.Tone{
			Fundamental: mlsic.Sine{
				Frequency: freq,
				Amplitude: 0.4 * i,
				Duration:  time.Duration(duration) * time.Millisecond,
			},
			Panning: pan,
		}
	}

	for i := 1.; i > 0.; i -= 0.1 {
		toneIndex += voice[toneIndex].Fundamental.DurationInSamples()

		voice[toneIndex] = mlsic.Tone{
			Fundamental: mlsic.Sine{
				Frequency: freq,
				Amplitude: 0.4 * i,
				Duration:  time.Duration(5 * time.Millisecond),
			},
			Panning: pan,
		}
	}

	return toneIndex
}

func move3(toneIndex int, voices ...mlsic.Voice) []mlsic.Voice {
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

func move3UpDown(freq float64, pan float64, tone int, voice mlsic.Voice, factor1, factor2 float64, duration int) int {
	for i := 0.; i < 1.; i += factor1 {
		tone += voice[tone].Fundamental.DurationInSamples()

		voice[tone] = mlsic.Tone{
			Fundamental: mlsic.Sine{
				Frequency: freq,
				Amplitude: i / 8,
				Duration:  time.Duration(duration) * time.Millisecond,
			},
			Panning: pan,
		}
	}

	for i := 1.; i > 0.; i -= factor2 {
		tone += voice[tone].Fundamental.DurationInSamples()

		voice[tone] = mlsic.Tone{
			Fundamental: mlsic.Sine{
				Frequency: freq,
				Amplitude: i / 8,
				Duration:  time.Duration(5 * time.Millisecond),
			},
			Panning: pan,
		}
	}

	return tone
}

func move4(toneIndex int, voices ...mlsic.Voice) []mlsic.Voice {
	for _, voice := range voices {
		toneIndex += voice[toneIndex].Fundamental.DurationInSamples()

		voice[toneIndex] = mlsic.Tone{
			Fundamental: mlsic.Sine{
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

		voice[toneIndex] = mlsic.Tone{
			Fundamental: mlsic.Sine{
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

		voice[toneIndex] = mlsic.Tone{
			Fundamental: mlsic.Sine{
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

		voice[toneIndex] = mlsic.Tone{
			Fundamental: mlsic.Sine{
				Frequency: freq,
				Amplitude: .2 / 8,
				Duration:  time.Duration(1000) * time.Millisecond,
			},
			Panning: pan,
		}
	}

	return voices
}
