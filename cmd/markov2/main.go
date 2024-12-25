package main

import (
	"flag"
	"math/rand"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/bh90210/mlsic"
	"github.com/bh90210/mlsic/markov"
	"github.com/bh90210/mlsic/markov/seed"
	"github.com/bh90210/mlsic/render"
	"github.com/mb-14/gomarkov"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	debug := flag.Bool("debug", false, "sets log level to debug")
	// ngenerations := flag.Int("ngen", 2, "sets log level to debug")
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

	freqModelPath := path.Join(*seedModelPath, "freq.json")
	ampModelPath := path.Join(*seedModelPath, "amp.json")
	durModelPath := path.Join(*seedModelPath, "dur.json")
	panModelPath := path.Join(*seedModelPath, "pan.json")

	freqModel, err := os.ReadFile(freqModelPath)
	if err != nil {
		log.Fatal().Interface("seedModel", freqModel).Err(err).Msg("reading seed")
	}

	ampModel, err := os.ReadFile(ampModelPath)
	if err != nil {
		log.Fatal().Interface("seedModel", freqModel).Err(err).Msg("reading seed")
	}

	durModel, err := os.ReadFile(durModelPath)
	if err != nil {
		log.Fatal().Interface("seedModel", freqModel).Err(err).Msg("reading seed")
	}

	panModel, err := os.ReadFile(panModelPath)
	if err != nil {
		log.Fatal().Interface("seedModel", freqModel).Err(err).Msg("reading seed")
	}

	freq := gomarkov.NewChain(1)
	err = freq.UnmarshalJSON(freqModel)
	if err != nil {
		log.Fatal().Err(err).Msg("unmarshaling seed")
	}

	amp := gomarkov.NewChain(1)
	err = amp.UnmarshalJSON(ampModel)
	if err != nil {
		log.Fatal().Err(err).Msg("unmarshaling seed")
	}

	dur := gomarkov.NewChain(1)
	err = dur.UnmarshalJSON(durModel)
	if err != nil {
		log.Fatal().Err(err).Msg("unmarshaling seed")
	}

	pan := gomarkov.NewChain(1)
	err = pan.UnmarshalJSON(panModel)
	if err != nil {
		log.Fatal().Err(err).Msg("unmarshaling seed")
	}

	noOfVoices := 4
	voices := make(map[int]markov.Voice, noOfVoices)
	for i := 0; i < noOfVoices; i++ {
		voices[i] = make(markov.Voice)
	}

	ngram := []string{"280.000000", "0.000000", "2200", "0.500000"}
	for u := 0; u < noOfVoices; u++ {
		var toneIndex int
		for o := 0; o < 100; o++ {
			for i := 0; i < 100; i++ {

				values := nextValueGen(freq, amp, dur, pan, ngram, int64(o))

				freq, err := strconv.ParseFloat(values[0], 64)
				if err != nil {
					log.Fatal().Err(err).Msg("parsing frequency")
				}

				amp, err := strconv.ParseFloat(values[1], 64)
				if err != nil {
					log.Fatal().Err(err).Msg("parsing amplitude")
				}

				dur, err := strconv.Atoi(values[2])
				if err != nil {
					log.Fatal().Err(err).Msg("parsing duration")
				}

				pan, err := strconv.ParseFloat(values[3], 64)
				if err != nil {
					log.Fatal().Err(err).Msg("parsing panning")
				}

				voices[u][toneIndex] = markov.Tone{
					Fundamental: markov.Sine{
						Frequency: freq,
						Amplitude: amp,
						Duration:  time.Duration(dur/mlsic.SignalLengthMultiplier) * time.Millisecond,
					},
					Panning: pan,
				}

				toneIndex += voices[u][toneIndex].Fundamental.DurationInSamples()

				ngram = values
			}
		}
	}

	for key, voice := range voices {
		voices[key] = seed.Partials(voice, seed.PrimeMove1)
	}

	var poly []markov.Voice
	poly = append(poly, voices[0], voices[1], voices[2], voices[3])

	// Export the new model.
	m := &markov.Model{
		Freq: freq,
		Dur:  dur,
		Amp:  amp,
		Pan:  pan,
	}

	m.Add2(poly)

	// Save seed model.
	err = m.Export(*modelsPath)
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

func nextValueGen(freq, amp, dur, pan *gomarkov.Chain, ngram []string, seed int64) []string {
	r := rand.New(rand.NewSource(seed))

	freqOut, err := freq.GenerateDeterministic([]string{ngram[0]}, r)
	if err != nil {
		log.Fatal().Err(err).Int64("seed", seed).Msg("generating freq")
	}

	if freqOut == "$" {
		freqOut = ngram[0]
	}

	ampOut, err := amp.GenerateDeterministic([]string{ngram[1]}, r)
	if err != nil {
		log.Fatal().Err(err).Int64("seed", seed).Msg("generating amp")
	}

	if ampOut == "$" {
		ampOut = ngram[1]
	}

	durOut, err := dur.GenerateDeterministic([]string{ngram[2]}, r)
	if err != nil {
		log.Fatal().Err(err).Int64("seed", seed).Msg("generating dur")
	}

	if durOut == "$" {
		durOut = ngram[2]
	}

	panOut, err := pan.GenerateDeterministic([]string{ngram[3]}, r)
	if err != nil {
		log.Fatal().Err(err).Int64("seed", seed).Msg("generating pan")
	}

	if panOut == "$" {
		panOut = ngram[3]
	}

	return []string{freqOut, ampOut, durOut, panOut}
}
