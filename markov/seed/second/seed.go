package main

import (
	"time"

	"github.com/bh90210/mlsic/markov"
	"github.com/bh90210/mlsic/markov/seed"
	"github.com/rs/zerolog/log"
)

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
	voice1 = seed.Partials(voice1, seed.PrimeMove1)
	voice2 = seed.Partials(voice2, seed.PrimeMove1)
	voice3 = seed.Partials(voice3, seed.PrimeMove1)
	voice4 = seed.Partials(voice4, seed.PrimeMove1)

	// Append voices to poly slice.
	poly = append(poly, voice1, voice2, voice3, voice4)

	return poly
}

func injectStart(index int, voice markov.Voice) int {
	voice[index] = markov.Tone{
		Fundamental: markov.Sine{
			Frequency: 180.000000,
			Amplitude: 0.100000,
			Duration:  time.Duration(50) * time.Millisecond,
		},
		Panning: 0.100000,
	}

	index += voice[index].Fundamental.DurationInSamples()

	return index
}

func upDown(freq float64, pan float64, toneIndex int, voice markov.Voice, factor float64, duration int) int {
	for i := 0.; i < 1.; i += factor {
		toneIndex += voice[toneIndex].Fundamental.DurationInSamples()

		voice[toneIndex] = markov.Tone{
			Fundamental: markov.Sine{
				Frequency: freq,
				Amplitude: 0.4 * i,
				Duration:  time.Duration(duration) * time.Millisecond,
			},
			Panning: pan,
		}
	}

	for i := 1.; i > 0.; i -= 0.1 {
		toneIndex += voice[toneIndex].Fundamental.DurationInSamples()

		voice[toneIndex] = markov.Tone{
			Fundamental: markov.Sine{
				Frequency: freq,
				Amplitude: 0.4 * i,
				Duration:  time.Duration(5 * time.Millisecond),
			},
			Panning: pan,
		}
	}

	return toneIndex
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
