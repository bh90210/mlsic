// Package pan holds various mlsic.Pan implementations.
package pan

import (
	"testing"

	"github.com/bh90210/mlsic/v1"
	"github.com/go-audio/audio"
	"github.com/stretchr/testify/require"
)

var tests = map[string]struct {
	pan          mlsic.Pan
	signal       audio.PCMBuffer
	position     float32
	wantSignal   []*audio.PCMBuffer
	wantChannels int
}{
	"mono": {
		pan:          &Mono{},
		signal:       audio.PCMBuffer{F32: []float32{1}},
		position:     0,
		wantSignal:   []*audio.PCMBuffer{{F32: []float32{1}}},
		wantChannels: 1,
	},
	"stereo linear (right)": {
		pan:          &LinearStereo{},
		signal:       audio.PCMBuffer{F32: []float32{1}},
		position:     1,
		wantSignal:   []*audio.PCMBuffer{{F32: []float32{0}}, {F32: []float32{1}}},
		wantChannels: 2,
	},
	"stereo linear (middle)": {
		pan:          &LinearStereo{},
		signal:       audio.PCMBuffer{F32: []float32{1}},
		position:     0.5,
		wantSignal:   []*audio.PCMBuffer{{F32: []float32{0.5}}, {F32: []float32{0.5}}},
		wantChannels: 2,
	},
}

func TestApply(t *testing.T) {
	r := require.New(t)

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.pan.Apply(tc.signal, tc.position)

			r.Len(got, tc.wantChannels)
			r.Equal(tc.wantSignal, got)
		})
	}
}

func TestChannels(t *testing.T) {
	r := require.New(t)

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.pan.Channels()

			r.Equal(tc.wantChannels, got)
		})
	}
}
