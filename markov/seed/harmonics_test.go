package seed

import (
	"testing"
	"time"

	"github.com/bh90210/mlsic/markov"
	"github.com/stretchr/testify/assert"
)

func TestPartials(t *testing.T) {
	var poly []markov.Voice

	voice := make(markov.Voice)
	voice[0] = markov.Tone{
		Fundamental: markov.Sine{
			Frequency: 440,
			Amplitude: 1.,
			Duration:  time.Duration(1 * time.Millisecond),
		},
		Panning: 0.5,
	}

	poly = append(poly, voice)

	h := PrimeHarmonics{}
	got := h.Partials(poly)

	assert.Equal(t, []markov.Voice{}, got)
}
