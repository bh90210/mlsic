package seed

import (
	"testing"

	"github.com/bh90210/mlsic"
	"github.com/bh90210/mlsic/markov"
	"github.com/stretchr/testify/assert"
)

type testSeedStruct struct {
	trains markov.Trains
	left   mlsic.Audio
	right  mlsic.Audio
}

func TestDeconstructTrains(t *testing.T) {
	tests := []testSeedStruct{
		{
			trains: MelodyTrain(),
			left:   []float64{0},
			right:  []float64{1},
		},
	}

	for _, test := range tests {
		left, right := DeconstructTrains(test.trains)
		assert.Equal(t, test.left, left)
		assert.Equal(t, test.right, right)
	}
}
