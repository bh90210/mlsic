package main

import (
	"math/rand"
	"os"

	"github.com/mb-14/gomarkov"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	durModel, err := os.ReadFile("/media/disk/mlsic/markov/seed/second/dur.json")
	if err != nil {
		log.Fatal().Interface("seedModel", durModel).Err(err).Msg("reading seed")
	}

	m := gomarkov.NewChain(1)
	err = m.UnmarshalJSON(durModel)
	if err != nil {
		log.Fatal().Err(err).Msg("unmarshaling seed")
	}

	found := map[string]map[string]int{}
	ngram := []string{"2200"}
	for i := 0; i < 1000; i++ {
		r := rand.New(rand.NewSource(int64(i)))
		o, e := m.GenerateDeterministic(ngram, r)
		if e != nil {
			log.Fatal().Err(e).Msg("generating seed")
		}

		if found[o] == nil {
			found[o] = map[string]int{}
		}

		// Add the output to the found map, including the ngram that produced it.
		found[o][ngram[0]]++

		ngram = append(ngram[1:], o)
	}

	for k, v := range found {
		log.Info().Str("key", k).Interface("v", v).Int("len", len(v)).Msg("v")
	}
}
