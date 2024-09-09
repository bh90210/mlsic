package markov

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"sync"
	"time"

	"github.com/bh90210/mlsic"
	"github.com/mb-14/gomarkov"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Song works as the entrance to the Markov chains generation.
type Song struct {
	// NGenerations is the number of generations to generate based on the initial seed.
	NGenerations int
	// FilePath is the directory the audio files wil be saved.
	FilePath string
	// ModelsPath is the directory the models produced out of each ngeneration wil be saved.
	ModelsPath string
	// SeedModelPath is the path of the initial seed jsons.
	SeedModelPath string

	// Harmonics is the harmonics structure that will be used for audio generation.
	Harmonics *Harmonics
}

type model struct {
	Int      int `json:"int"`
	SpoolMap any `json:"spool_map"`
	FreqMat  any `json:"freq_mat"`
}

// NGen will process the seed model and based on it will generate the appropriate amount of generation cycles.
func (s *Song) NGen() {
	log.Info().Msg("NGen")

	// Generate a new model and audio output for each generation.
	for i := 0; i < s.NGenerations; i++ {
		log.Logger = log.With().Int("gen", i).Logger()

		log.Info().Msg("NGen")

		var freqModel string
		var ampModel string
		var durModel string

		index := strconv.Itoa(i)

		// If we are on the first iteration we must start by reading
		// the seed model.
		if i == 0 {
			freqModel = filepath.Join(s.SeedModelPath, "freq.json")
			ampModel = filepath.Join(s.SeedModelPath, "amp.json")
			durModel = filepath.Join(s.SeedModelPath, "dur.json")

			// If the seed model is already processed read the previously generated model.
		} else {
			previousModelsIndex := strconv.Itoa(i - 1)

			freqModel = filepath.Join(s.ModelsPath, "gen"+previousModelsIndex, "freq.json")
			ampModel = filepath.Join(s.ModelsPath, "gen"+previousModelsIndex, "amp.json")
			durModel = filepath.Join(s.ModelsPath, "gen"+previousModelsIndex, "dur.json")
		}

		log.Info().Msg("reading files")

		freq, err := os.ReadFile(freqModel)
		if err != nil {
			log.Fatal().Err(err).Msg("reading freq")
		}

		amp, err := os.ReadFile(ampModel)
		if err != nil {
			log.Fatal().Err(err).Msg("reading amp")
		}

		dur, err := os.ReadFile(durModel)
		if err != nil {
			log.Fatal().Err(err).Msg("reading dur")
		}

		log.Info().Msg("creating chains")

		// Prepare a Markov Train so we can load the previously created model.
		t := Models{
			Freq: gomarkov.NewChain(1),
			Amp:  gomarkov.NewChain(1),
			Dur:  gomarkov.NewChain(1),
		}

		// Load previously generated model.
		t.Freq.UnmarshalJSON(freq)
		t.Amp.UnmarshalJSON(amp)
		t.Dur.UnmarshalJSON(dur)

		var wg sync.WaitGroup

		var generationFreqs [][]float64
		var generationAmps [][]float64
		var generationDurs [][]float64

		for i := 0; i < 3; i++ {
			wg.Add(1)

			go func(i int) {
				defer wg.Done()

				switch i {
				case 0:
					l := log.Logger
					l = l.With().Str("field", "freq").Logger()

					// Generate new values for frequencies, amplitudes and durations based on previous model.
					var frequencies model
					err = json.Unmarshal(freq, &frequencies)
					if err != nil {
						l.Fatal().Err(err).Msg("unmarsh freq")
					}

					l.Info().Msg("entering loop")

					generationFreqs, err = markovGenerator(l, frequencies.SpoolMap, t.Freq)
					if err != nil {
						l.Fatal().Err(err).Msg("freq loop")
					}

				case 1:
					l := log.Logger
					l = l.With().Str("field", "amp").Logger()

					var amplitudes model
					err = json.Unmarshal(amp, &amplitudes)
					if err != nil {
						l.Fatal().Err(err).Msg("unmarsh amp")
					}

					l.Info().Msg("entering loop")

					generationAmps, err = markovGenerator(l, amplitudes.SpoolMap, t.Amp)
					if err != nil {
						l.Fatal().Err(err).Msg("amp loop")
					}

				case 2:
					l := log.Logger
					l = l.With().Str("field", "dur").Logger()

					var durations model
					err = json.Unmarshal(dur, &durations)
					if err != nil {
						l.Fatal().Err(err).Msg("unmarsh dur")
					}

					l.Info().Msg("entering loop")

					generationDurs, err = markovGenerator(l, durations.SpoolMap, t.Dur, true)
					if err != nil {
						l.Fatal().Err(err).Msg("dur loop")
					}

				}
			}(i)
		}

		wg.Wait()

		// Reset logger to remove "field".
		log.Logger = log.With().Reset().Logger().With().Int("gen", i).Logger()

		log.Info().Msg("creating sines train")

		// Create sines train.
		var train []mlsic.Sine

		for i, freqs := range generationFreqs {
			var outOfBoundsAmp bool
			var outOfBoundsDur bool

			if len(generationAmps)-1 < i {
				outOfBoundsAmp = true
			}

			if len(generationDurs)-1 < i {
				outOfBoundsDur = true
			}

			for o, freq := range freqs {
				// TODO: 0 is arbitrary, fix it.
				amp := 0.
				if !outOfBoundsAmp && !(len(generationAmps[i])-1 < o) {
					amp = generationAmps[i][o]
				}

				// TODO: 10 is arbitrary, fix it.
				dur := 10.
				if !outOfBoundsDur && !(len(generationDurs[i])-1 < o) {
					dur = generationDurs[i][o]
				}

				train = append(train, mlsic.Sine{
					Frequency: freq,
					Amplitude: amp,
					Duration:  time.Duration(dur) * time.Millisecond,
				})
			}
		}

		log.Info().Msg("audio files gen")

		// filePath := filepath.Join(s.FilePath, "gen"+index)

		err = os.MkdirAll(s.FilePath, 0755)
		if err != nil {
			log.Fatal().Err(err).Msg("creating audio directory")
		}

		// Generate audio based on the new model.
		Generate(s.FilePath, train, s.Harmonics, i)

		log.Info().Msg("export models")

		// Save the new model.
		t.Add(train)

		modelsPath := filepath.Join(s.ModelsPath, "gen"+index)

		err = os.MkdirAll(modelsPath, 0755)
		if err != nil {
			log.Fatal().Err(err).Msg("creating models directory")
		}

		err = t.Export(modelsPath)
		if err != nil {
			log.Fatal().Err(err).Msg("exporting seed")
		}
	}
}

// TODO: better name.
func markovGenerator(l zerolog.Logger, spoolMap any, chain *gomarkov.Chain, dur ...bool) ([][]float64, error) {
	mapped := spoolMap.(map[string]interface{})

	// Sort mapped.
	var sortedMapped []float64
	for k := range mapped {
		if k == "$" || k == "^" || k == "+Inf" {
			continue
		}

		flo, err := strconv.ParseFloat(k, 64)
		if err != nil {
			return nil, fmt.Errorf("parse float %w", err)
		}

		sortedMapped = append(sortedMapped, flo)
	}

	slices.Sort(sortedMapped)

	var mu sync.Mutex
	var wg sync.WaitGroup

	temporaryTrain := make([][]float64, len(sortedMapped))
	wg.Add(len(sortedMapped))

	for o, value := range sortedMapped {
		go func(o int, value float64) {
			defer wg.Done()

			starting := []string{fmt.Sprintf("%f", value)}

			if len(dur) > 0 {
				starting = []string{fmt.Sprintf("%.0f", value)}
			}

			var temp []float64
			for i := 0; ; i++ {
				l := l.With().
					Int("outer iter", o).
					Int("inner iter", i).
					Logger()

				// generated, err := chain.GenerateDeterministic(starting, rand.New(rand.NewSource(int64(o))))
				generated, err := chain.GenerateDeterministic(starting, rand.New(rand.NewSource(int64(420))))
				if err != nil {
					l.Fatal().Err(err).Msg("generating next markov")
				}

				if generated == "$" {
					l.Debug().Msg("$")
					break
				}

				flo, err := strconv.ParseFloat(generated, 64)
				if err != nil {
					l.Fatal().Err(err).Msg("parsing string to float")
				}

				// Check is we are looping.
				foundPattern := patternFinder(temp, flo)
				if foundPattern {
					l.Debug().
						Str("generated", generated).
						Msg("loop found")
					break
				}

				temp = append(temp, flo)

				starting = []string{generated}
			}

			mu.Lock()
			temporaryTrain[o] = temp
			mu.Unlock()

		}(o, value)
	}

	wg.Wait()

	return temporaryTrain, nil
}

func patternFinder(temp []float64, flo float64) bool {
	temp = append(temp, flo)

	if len(temp) < 10 {
		return false
	}

	var trackPattern int

	found := make(map[int]bool)
	for i, v := range temp {
		for ii, vv := range temp {
			if i == ii {
				continue
			}

			_, ok := found[ii]
			if ok {
				continue
			}

			if v == vv {
				found[ii] = true
			}
		}
	}

	var foundSorted []int
	for k := range found {
		foundSorted = append(foundSorted, k)
	}

	slices.Sort(foundSorted)

	// i = the length of the slice. When we are at index [1] stop.
	for i := len(foundSorted) - 1; i > 0; i-- {
		if foundSorted[i]-foundSorted[i-1] == 1 {
			trackPattern++
		}
	}

	if trackPattern == len(foundSorted)-1 {
		return true
	}

	return false
}
